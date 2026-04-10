package mapper

import (
	"fmt"

	"github.com/lucas/netmap/internal/classifier"
	"github.com/lucas/netmap/internal/models"
)

// Mapper Engine handles constructing the map graph.
type Mapper struct {
	Graph models.MapGraph
}

func NewMapper(domain string) *Mapper {
	return &Mapper{
		Graph: models.MapGraph{
			Target: models.DomainTarget{Domain: domain},
			Nodes:  []models.Node{},
			Edges:  []models.Edge{},
		},
	}
}

// Run executes the mapping process. Currently generates stub data for V1.
func (m *Mapper) Run() {
	// Add Root Node
	rootID := "root-" + m.Graph.Target.Domain
	m.addNode(models.Node{
		ID:    rootID,
		Label: m.Graph.Target.Domain,
		Type:  models.RootNode,
	})

	// V1 Stub: Generate mock subdomains
	subDoms := []string{"api", "admin", "www"}
	for _, sub := range subDoms {
		subName := fmt.Sprintf("%s.%s", sub, m.Graph.Target.Domain)
		subID := "sub-" + subName
		m.addNode(models.Node{
			ID:       subID,
			Label:    subName,
			Type:     models.SubdomainNode,
			ParentID: rootID,
		})
		m.addEdge(rootID, subID, "has_subdomain")

		// Mock Endpoints per subdomain
		var endpoints []string
		if sub == "api" {
			endpoints = []string{"/v1/users", "/auth/login"}
		} else if sub == "admin" {
			endpoints = []string{"/login"}
		} else if sub == "www" {
			endpoints = []string{"/", "/dashboard", "/about"}
		}

		for _, ep := range endpoints {
			epID := fmt.Sprintf("ep-%s-%s", subName, ep)
			m.addNode(models.Node{
				ID:       epID,
				Label:    ep,
				Type:     models.EndpointNode,
				Category: classifier.ClassifyEndpoint(ep),
				ParentID: subID,
			})
			m.addEdge(subID, epID, "has_endpoint")
		}
	}
}

func (m *Mapper) addNode(n models.Node) {
	m.Graph.Nodes = append(m.Graph.Nodes, n)
}

func (m *Mapper) addEdge(source, target, rel string) {
	m.Graph.Edges = append(m.Graph.Edges, models.Edge{
		Source:           source,
		Target:           target,
		RelationshipType: rel,
	})
}
