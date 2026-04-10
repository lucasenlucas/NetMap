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
	Graph          models.MapGraph
	Mode           models.FocusMode
	PackName       string
	CustomWordlist string
	mu             sync.Mutex
	client         *http.Client
}

func NewMapper(domain string, mode string, pack string, customWordlist string) *Mapper {
	return &Mapper{
		Graph: models.MapGraph{
			Target: models.DomainTarget{Domain: domain},
			Nodes:  []models.Node{},
			Edges:  []models.Edge{},
		},
		Mode:           models.FocusMode(mode),
		PackName:       pack,
		CustomWordlist: customWordlist,
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
	if m.PackName != "" && m.PackName != "standard" {
		fmt.Fprintf(os.Stderr, "[~] Active Pack: %s\n", m.PackName)
	}

	// 2. Resolve Discovery Lists
	subList, pathList := discovery.GetPack(m.PackName)
	if m.CustomWordlist != "" {
		fmt.Fprintf(os.Stderr, "[~] Loading custom wordlist: %s\n", m.CustomWordlist)
		custom, err := discovery.LoadWordlistFromFile(m.CustomWordlist)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error loading wordlist: %v\n", err)
		} else {
			pathList = append(pathList, custom...)
		}
	}

	// 3. Discover Subdomains
	subdomains := m.discoverSubdomains(subList)
	subdomains = append([]string{m.Graph.Target.Domain}, subdomains...)

	fmt.Fprintf(os.Stderr, "[~] Discovered %d unique hosts. Commencing active service mapping...\n", len(subdomains))

	// 4. Multi-threaded Processing with Concurrency Limit
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 50)

	for _, sub := range subdomains {
		wg.Add(1)
		go func(subdomain string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

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

			// Endpoint enumeration
			var epWg sync.WaitGroup
			epSemaphore := make(chan struct{}, 10)
			for _, ep := range pathList {
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

func (m *Mapper) discoverSubdomains(wordlist []string) []string {
	uniqueSubs := make(map[string]bool)
	var finalSubs []string

	// OSINT Phase
	fmt.Fprintf(os.Stderr, "[+] Gathering OSINT intelligence (crt.sh)...\n")
	for _, s := range m.fetchOSINTSubdomains(m.Graph.Target.Domain) {
		if !uniqueSubs[s] {
			uniqueSubs[s] = true
			finalSubs = append(finalSubs, s)
		}
	}

	// DNS Phase (if advanced OR specific pack that needs it)
	if m.Mode == "advanced" || m.PackName == "dns-extended" {
		fmt.Fprintf(os.Stderr, "[+] Triggering DNS brute-force discovery...\n")
		var dnsWg sync.WaitGroup
		dnsMu := sync.Mutex{}
		for _, prefix := range wordlist {
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

	limit := 20
	if m.Mode == "advanced" {
		limit = 150
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
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
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
