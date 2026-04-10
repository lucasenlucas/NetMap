package models

// NodeType classifies the level of the node in the map.
type NodeType string

const (
	RootNode      NodeType = "Root"
	SubdomainNode NodeType = "Subdomain"
	EndpointNode  NodeType = "Endpoint"
)

// FocusMode represents the filter applied to the mapping output.
type FocusMode string

const (
	FocusAll    FocusMode = "all"
	FocusAuth   FocusMode = "auth"
	FocusAdmin  FocusMode = "admin"
	FocusAPI    FocusMode = "api"
	FocusConfig FocusMode = "config"
	FocusDev    FocusMode = "dev"
)

// EndpointType classifies the nature of an endpoint.
type EndpointType string

const (
	TypeAuth    EndpointType = "auth"
	TypeAdmin   EndpointType = "admin"
	TypeAPI     EndpointType = "api"
	TypeConfig  EndpointType = "config"
	TypeDev     EndpointType = "dev"
	TypeGeneral EndpointType = "general"
	TypeUnknown EndpointType = "unknown"
)

// Node represents a single entity in the parsed network structure (domain, subdomain, endpoint).
type Node struct {
	ID        string       `json:"id"`
	Label     string       `json:"label"`
	Type      NodeType     `json:"type"`
	Category  EndpointType `json:"category,omitempty"`
	ParentID  string       `json:"parent_id,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	RiskLevel int          `json:"risk_level,omitempty"`
}

// Edge represents a relationship between two nodes.
type Edge struct {
	Source           string `json:"source"`
	Target           string `json:"target"`
	RelationshipType string `json:"relationship_type"`
}

// MapGraph is the core structure holding the mapped data.
type MapGraph struct {
	Target DomainTarget `json:"target"`
	Nodes  []Node       `json:"nodes"`
	Edges  []Edge       `json:"edges"`
}

// DomainTarget holds information about the root scan target.
type DomainTarget struct {
	Domain string `json:"domain"`
}

// Endpoint represents a scanned endpoint.
type Endpoint struct {
	Path     string       `json:"path"`
	Category EndpointType `json:"category"`
}
