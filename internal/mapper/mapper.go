package mapper

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lucas/netmap/internal/classifier"
	"github.com/lucas/netmap/internal/discovery"
	"github.com/lucas/netmap/internal/models"
)

// CrtShResponse represents a single valid entry block from crt.sh
type CrtShResponse struct {
	NameValue string `json:"name_value"`
}

// Mapper Engine handles constructing the live map graph.
type Mapper struct {
	Graph  models.MapGraph
	Mode   models.FocusMode // We use FocusMode or a dedicated Mode type; Let's reuse models or define basic/advanced
	mu     sync.Mutex
	client *http.Client
}

func NewMapper(domain string, mode string) *Mapper {
	return &Mapper{
		Graph: models.MapGraph{
			Target: models.DomainTarget{Domain: domain},
			Nodes:  []models.Node{},
			Edges:  []models.Edge{},
		},
		Mode: models.FocusMode(mode),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Run executes the live mapping process utilizing OSINT and Probing.
func (m *Mapper) Run() {
	// 1. Add Root Node
	rootID := "root-" + m.Graph.Target.Domain
	m.addNode(models.Node{
		ID:    rootID,
		Label: m.Graph.Target.Domain,
		Type:  models.RootNode,
	})

	fmt.Fprintf(os.Stderr, "[~] Initializing discovery for %s...\n", m.Graph.Target.Domain)

	// 2. Discover Subdomains
	subdomains := m.discoverSubdomains()
	
	// Ensure root domain is always checked
	subdomains = append([]string{m.Graph.Target.Domain}, subdomains...)

	fmt.Fprintf(os.Stderr, "[~] Discovered %d unique subdomains. Commencing active probing...\n", len(subdomains)-1)

	// 3. Define Endpoints to Check
	endpointsToCheck := []string{"/", "/admin", "/login", "/dashboard", "/api", "/api/v1", "/graphql", "/auth"}
	if m.Mode == "advanced" {
		endpointsToCheck = discovery.GetPaths()
	}

	// 4. Multi-threaded Processing with Concurrency Limit
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50) // Max 50 concurrent subdomain tasks

	for _, sub := range subdomains {
		wg.Add(1)
		go func(subdomain string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Quick alive check
			if !m.isAlive(subdomain) {
				return
			}

			subID := "sub-" + subdomain
			if subdomain != m.Graph.Target.Domain {
				m.addNodeSafe(models.Node{
					ID:       subID,
					Label:    subdomain,
					Type:     models.SubdomainNode,
					ParentID: rootID,
				})
				m.addEdgeSafe(rootID, subID, "has_subdomain")
			} else {
				subID = rootID
			}

			// Endpoint enumeration with sub-semaphore
			var epWg sync.WaitGroup
			epSemaphore := make(chan struct{}, 10) // Max 10 concurrent endpoint probes per subdomain

			for _, ep := range endpointsToCheck {
				epWg.Add(1)
				go func(endpoint string) {
					defer epWg.Done()
					epSemaphore <- struct{}{}
					defer func() { <-epSemaphore }()

					if m.probeEndpoint(subdomain, endpoint) {
						epID := fmt.Sprintf("ep-%s-%s", subdomain, endpoint)
						m.addNodeSafe(models.Node{
							ID:       epID,
							Label:    endpoint,
							Type:     models.EndpointNode,
							Category: classifier.ClassifyEndpoint(endpoint),
							ParentID: subID,
						})
						m.addEdgeSafe(subID, epID, "has_endpoint")
					}
				}(ep)
			}
			epWg.Wait()
		}(sub)
	}

	wg.Wait()
}

func (m *Mapper) discoverSubdomains() []string {
	uniqueSubs := make(map[string]bool)
	var finalSubs []string

	// Phase 1: OSINT (Always)
	fmt.Fprintf(os.Stderr, "[+] Querying OSINT intelligence (crt.sh)...\n")
	osintSubs := m.fetchOSINTSubdomains(m.Graph.Target.Domain)
	for _, s := range osintSubs {
		if !uniqueSubs[s] {
			uniqueSubs[s] = true
			finalSubs = append(finalSubs, s)
		}
	}

	// Phase 2: DNS Brute Force (Only in Advanced Mode)
	if m.Mode == "advanced" {
		fmt.Fprintf(os.Stderr, "[+] Commencing DNS discovery for top 50 subdomains...\n")
		var dnsWg sync.WaitGroup
		dnsMu := sync.Mutex{}
		
		for _, prefix := range discovery.GetSubdomains() {
			dnsWg.Add(1)
			go func(p string) {
				defer dnsWg.Done()
				sub := fmt.Sprintf("%s.%s", p, m.Graph.Target.Domain)
				if _, err := net.LookupHost(sub); err == nil {
					dnsMu.Lock()
					if !uniqueSubs[sub] {
						uniqueSubs[sub] = true
						finalSubs = append(finalSubs, sub)
					}
					dnsMu.Unlock()
				}
			}(prefix)
		}
		dnsWg.Wait()
	}

	// Cap for stability if not in advanced mode, or higher cap in advanced
	limit := 15
	if m.Mode == "advanced" {
		limit = 100
	}
	if len(finalSubs) > limit {
		finalSubs = finalSubs[:limit]
	}

	return finalSubs
}

func (m *Mapper) fetchOSINTSubdomains(domain string) []string {
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	resp, err := m.client.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return []string{}
	}
	defer resp.Body.Close()

	var results []CrtShResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return []string{}
	}

	var subs []string
	for _, r := range results {
		sub := strings.ReplaceAll(r.NameValue, "*.", "")
		sub = strings.TrimSpace(sub)
		if strings.Contains(sub, "\n") {
			sub = strings.Split(sub, "\n")[0]
		}
		if sub != "" && sub != domain && !strings.Contains(sub, "@") && !strings.Contains(sub, " ") {
			subs = append(subs, sub)
		}
	}
	return subs
}

// isAlive checks port 443 then 80
func (m *Mapper) isAlive(host string) bool {
	resp, err := m.client.Head("https://" + host)
	if err == nil {
		resp.Body.Close()
		return true
	}
	resp, err = m.client.Head("http://" + host)
	if err == nil {
		resp.Body.Close()
		return true
	}
	return false
}

func (m *Mapper) probeEndpoint(host, path string) bool {
	urls := []string{"https://" + host + path, "http://" + host + path}
	for _, u := range urls {
		req, _ := http.NewRequest("GET", u, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; NetMap/1.0; +https://github.com/lucasenlucas)")
		resp, err := m.client.Do(req)
		if err != nil {
			continue
		}
		status := resp.StatusCode
		resp.Body.Close()
		if (status >= 200 && status < 300) || status == 401 || status == 403 {
			return true
		}
	}
	return false
}

func (m *Mapper) addNode(n models.Node) {
	m.Graph.Nodes = append(m.Graph.Nodes, n)
}

func (m *Mapper) addNodeSafe(n models.Node) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Graph.Nodes = append(m.Graph.Nodes, n)
}

func (m *Mapper) addEdgeSafe(source, target, rel string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Graph.Edges = append(m.Graph.Edges, models.Edge{
		Source:           source,
		Target:           target,
		RelationshipType: rel,
	})
}
