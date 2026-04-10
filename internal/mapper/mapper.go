package mapper

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
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
	
	// Progress Tracking
	totalSubs       int32
	completedSubs   int32
	totalPaths      int32
	completedPaths  int32
	activeHost      string
	
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

	fmt.Fprintf(os.Stderr, "[~] Initializing NetMap Intelligence for %s...\n", m.Graph.Target.Domain)

	// 2. Resolve Discovery Lists
	subList, pathList := discovery.GetPack(m.PackName)
	if m.CustomWordlist != "" {
		custom, err := discovery.LoadWordlistFromFile(m.CustomWordlist)
		if err == nil {
			pathList = append(pathList, custom...)
		}
	}

	// 3. Discover Subdomains
	subdomains := m.discoverSubdomains(subList)
	subdomains = append([]string{m.Graph.Target.Domain}, subdomains...)
	atomic.StoreInt32(&m.totalSubs, int32(len(subdomains)))
	atomic.StoreInt32(&m.totalPaths, int32(len(pathList)*len(subdomains)))

	// 4. Start Live Progress Reporter
	stopProgress := make(chan struct{})
	go m.progressReporter(stopProgress)

	// 5. Global Concurrency Management
	globalLimit := 100
	if m.PackName == "ultra" {
		globalLimit = 250 
	}
	globalSemaphore := make(chan struct{}, globalLimit)

	var wg sync.WaitGroup
	for _, sub := range subdomains {
		wg.Add(1)
		go func(subdomain string) {
			defer wg.Done()
			
			// Always lookup DNS records for discovered subdomains, even if HTTP is dead
			m.lookupDNSRecords(subdomain)

			if !m.isAlive(subdomain) {
				atomic.AddInt32(&m.completedSubs, 1)
				atomic.AddInt32(&m.completedPaths, int32(len(pathList)))
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
			for _, ep := range pathList {
				epWg.Add(1)
				globalSemaphore <- struct{}{}
				go func(endpoint string) {
					defer epWg.Done()
					defer func() { <-globalSemaphore }()
					m.setHost(subdomain)

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
					atomic.AddInt32(&m.completedPaths, 1)
				}(ep)
			}
			epWg.Wait()
			atomic.AddInt32(&m.completedSubs, 1)
		}(sub)
	}

	wg.Wait()
	close(stopProgress)
	time.Sleep(200 * time.Millisecond) 
	fmt.Fprintf(os.Stderr, "\r%s\r", strings.Repeat(" ", 80)) 
}

func (m *Mapper) lookupDNSRecords(host string) {
	parentID := "sub-" + host
	if host == m.Graph.Target.Domain {
		parentID = "root-" + host
	}

	// 1. Lookup CNAME
	if cname, err := net.LookupCNAME(host); err == nil && cname != "" && cname != host+"." {
		cnameID := fmt.Sprintf("dns-cname-%s", host)
		m.addNodeSafe(models.Node{
			ID:       cnameID,
			Label:    fmt.Sprintf("CNAME: %s", strings.TrimSuffix(cname, ".")),
			Type:     models.EndpointNode,
			Category: models.TypeDNS,
			ParentID: parentID,
		})
		m.addEdgeSafe(parentID, cnameID, "cname_points_to")
	}

	// 2. Lookup MX
	if mxs, err := net.LookupMX(host); err == nil {
		for _, mx := range mxs {
			mxID := fmt.Sprintf("dns-mx-%s-%s", host, mx.Host)
			m.addNodeSafe(models.Node{
				ID:       mxID,
				Label:    fmt.Sprintf("MX: %s (pref: %d)", strings.TrimSuffix(mx.Host, "."), mx.Pref),
				Type:     models.EndpointNode,
				Category: models.TypeDNS,
				ParentID: parentID,
			})
			m.addEdgeSafe(parentID, mxID, "mail_handled_by")
		}
	}

	// 3. Lookup TXT
	if txts, err := net.LookupTXT(host); err == nil {
		for i, txt := range txts {
			if len(txt) > 80 { txt = txt[:77] + "..." }
			txtID := fmt.Sprintf("dns-txt-%s-%d", host, i)
			m.addNodeSafe(models.Node{
				ID:       txtID,
				Label:    fmt.Sprintf("TXT: %s", txt),
				Type:     models.EndpointNode,
				Category: models.TypeDNS,
				ParentID: parentID,
			})
			m.addEdgeSafe(parentID, txtID, "has_txt_record")
		}
	}
}

func (m *Mapper) progressReporter(stop chan struct{}) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			subs := atomic.LoadInt32(&m.completedSubs)
			totalSubs := atomic.LoadInt32(&m.totalSubs)
			paths := atomic.LoadInt32(&m.completedPaths)
			totalPaths := atomic.LoadInt32(&m.totalPaths)
			host := m.getHost()

			fmt.Fprintf(os.Stderr, "\r\033[36m[~] Mapping: %d/%d hosts | %d/%d endpoints | [%s]\r\033[0m", 
				subs, totalSubs, paths, totalPaths, truncate(host, 15))
		}
	}
}

func (m *Mapper) discoverSubdomains(wordlist []string) []string {
	uniqueSubs := make(map[string]bool)
	var finalSubs []string

	fmt.Fprintf(os.Stderr, "[+] Gathering OSINT intelligence (crt.sh)...\n")
	for _, s := range m.fetchOSINTSubdomains(m.Graph.Target.Domain) {
		if !uniqueSubs[s] {
			uniqueSubs[s] = true
			finalSubs = append(finalSubs, s)
		}
	}

	if m.Mode == "advanced" || m.PackName == "dns-extended" || m.PackName == "full" || m.PackName == "ultra" {
		fmt.Fprintf(os.Stderr, "[+] Triggering DNS brute-force discovery (%d candidates)...\n", len(wordlist))
		var dnsWg sync.WaitGroup
		dnsMu := sync.Mutex{}
		dnsSemaphore := make(chan struct{}, 150)

		for _, prefix := range wordlist {
			dnsWg.Add(1)
			go func(p string) {
				defer dnsWg.Done()
				dnsSemaphore <- struct{}{}
				defer func() { <-dnsSemaphore }()

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

	limit := 30
	if m.PackName == "full" { limit = 150 }
	if m.PackName == "ultra" { limit = 500 }

	if len(finalSubs) > limit {
		finalSubs = finalSubs[:limit]
	}

	return finalSubs
}

func (m *Mapper) fetchOSINTSubdomains(domain string) []string {
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)
	resp, err := m.client.Get(url)
	if err != nil || resp.StatusCode != 200 { return []string{} }
	defer resp.Body.Close()

	var results []CrtShResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil { return []string{} }

	var subs []string
	for _, r := range results {
		sub := strings.ReplaceAll(r.NameValue, "*.", "")
		sub = strings.TrimSpace(sub)
		if strings.Contains(sub, "\n") { sub = strings.Split(sub, "\n")[0] }
		if sub != "" && sub != domain && !strings.Contains(sub, "@") && !strings.Contains(sub, " ") {
			subs = append(subs, sub)
		}
	}
	return subs
}

func (m *Mapper) isAlive(host string) bool {
	resp, err := m.client.Head("https://" + host)
	if err == nil { resp.Body.Close(); return true }
	resp, err = m.client.Head("http://" + host)
	if err == nil { resp.Body.Close(); return true }
	return false
}

func (m *Mapper) probeEndpoint(host, path string) bool {
	if !strings.HasPrefix(path, "/") { path = "/" + path }
	u := "https://" + host + path
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; NetMap/1.0; +https://github.com/lucasenlucas/NetMap)")
	resp, err := m.client.Do(req)
	if err != nil { return false }
	status := resp.StatusCode
	resp.Body.Close()
	return (status >= 200 && status < 300) || status == 401 || status == 403
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

func (m *Mapper) setHost(h string) {
	m.mu.Lock()
	m.activeHost = h
	m.mu.Unlock()
}

func (m *Mapper) getHost() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.activeHost
}

func truncate(s string, n int) string {
	if len(s) <= n { return s }
	return s[:n-3] + "..."
}
