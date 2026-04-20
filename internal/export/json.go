package export

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/netseries/netmap/internal/models"
)

// ExportWrapper encapsulates the graph payload in a flattened format ready for frontend consumption.
type ExportWrapper struct {
	Nodes    []models.Node  `json:"nodes"`
	Edges    []models.Edge  `json:"edges"`
	Metadata ExportMetadata `json:"metadata"`
}

type ExportMetadata struct {
	Domain     string `json:"domain"`
	TotalNodes int    `json:"totalNodes"`
	Generated  string `json:"generatedAt"`
	Version    string `json:"version"`
}

// PrintJSON outputs the map graph as a formatted JSON structure.
func PrintJSON(graph *models.MapGraph) {
	wrapper := ExportWrapper{
		Nodes: graph.Nodes,
		Edges: graph.Edges,
		Metadata: ExportMetadata{
			Domain:     graph.Target.Domain,
			TotalNodes: len(graph.Nodes),
			Generated:  time.Now().UTC().Format(time.RFC3339),
			Version:    "1.0.0",
		},
	}

	bytes, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		fmt.Printf(`{"error": "Failed to serialize JSON map: %v"}`, err)
		return
	}
	fmt.Println(string(bytes))
}
