package mapper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lucas/netmap/internal/classifier"
	"github.com/lucas/netmap/internal/models"
)

// CrtShResponse represents a single valid entry block from crt.sh
type CrtShResponse struct {
	NameValue string `json:"name_value"`
}

// Mapper Engine handles constructing the live map graph.
type Mapper struct {
	Graph  models.MapGraph
	mu     sync.Mutex
	client *http.Client
}

func NewMapper(domain string) *Mapper {
	return &Mapper{
		Graph: models.MapGraph{
			Target: models.DomainTarget{Domain: domain},
			Nodes:  []models.Node{},
			Edges:  []models.Edge{},
		},
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

	fmt.Fprintf(os.Stderr, "[~] Querying OSINT intelligence for %s...\n", m.Graph.Target.Domain)

	// 2. Discover Subdomains via CRT.sh
	subdomains := m.fetchSubdomains(m.Graph.Target.Domain)
	
	// Ensure root domain is always checked
	subdomains = append([]string{m.Graph.Target.Domain}, subdomains...)

	fmt.Fprintf(os.Stderr, "[~] Discovered %d subdomains. Commencing live ping sweep & endpoint probing...\n", len(subdomains)-1)

	var wg sync.WaitGroup
	
	endpointsToCheck := []string{
		"/",
		"/admin",
		"/login",
		"/dashboard",
		"/api",
		"/api/v1",
		"/graphql",
		"/auth",
	}

	for _, sub := range subdomains {
		wg.Add(1)
		go func(subdomain string) {
			defer wg.Done()
			
			// Quick root ping to verify host alive status
			if !m.isAlive(subdomain) {
				return // Keep output absolutely clean by discarding dead zones
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

			// Sub-thread endpoint enumeration
			var epWg sync.WaitGroup
			for _, ep := range endpointsToCheck {
				epWg.Add(1)
				go func(endpoint string) {
					defer epWg.Done()
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

func (m *Mapper) fetchSubdomains(domain string) []string {
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

	uniqueSubs := make(map[string]bool)
	var finalSubs []string

	for _, r := range results {
		sub := strings.ReplaceAll(r.NameValue, "*.", "")
		sub = strings.TrimSpace(sub)
		
		// Break out early if there's multiple lines in crt output (rare but possible mapping errors)
		if strings.Contains(sub, "\n") {
			parts := strings.Split(sub, "\n")
			sub = parts[0]
		}
		
		if sub != "" && sub != domain && !uniqueSubs[sub] && !strings.Contains(sub, "@") && !strings.Contains(sub, " ") {
			uniqueSubs[sub] = true
			finalSubs = append(finalSubs, sub)
		}
	}
	
	// Cap scale to 15 subdomains to ensure lightning fast terminal UX
	if len(finalSubs) > 15 {
		finalSubs = finalSubs[:15]
	}

	return finalSubs
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

// probeEndpoint verifies specific endpoint existence
func (m *Mapper) probeEndpoint(host, path string) bool {
	urls := []string{
		"https://" + host + path,
		"http://" + host + path,
	}

	for _, u := range urls {
		req, _ := http.NewRequest("GET", u, nil)
		// Spoof agent to bypass low-level scrap guards
		req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; NetMap/1.0; +https://github.com/lucasenlucas)")
		
		resp, err := m.client.Do(req)
		if err != nil {
			continue
		}
		
		status := resp.StatusCode
		resp.Body.Close()

		// Standard endpoints exist if 200 OK or Explicitly Gated (401, 403)
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
