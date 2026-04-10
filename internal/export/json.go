package export

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lucas/netmap/internal/models"
)

// ExportWrapper encapsulates the graph payload with metadata ensuring readiness for frontend consumption.
type ExportWrapper struct {
	Version   string           `json:"version"`
	Generated string           `json:"generated_at"`
	Target    string           `json:"target_domain"`
	Summary   ExportSummary    `json:"summary"`
	Graph     *models.MapGraph `json:"graph"`
}

type ExportSummary struct {
	TotalNodes int `json:"total_nodes"`
	TotalEdges int `json:"total_edges"`
}

// PrintJSON outputs the map graph as a formatted JSON structure.
func PrintJSON(graph *models.MapGraph) {
	wrapper := ExportWrapper{
		Version:   "1.0.0",
		Generated: time.Now().UTC().Format(time.RFC3339),
		Target:    graph.Target.Domain,
		Summary: ExportSummary{
			TotalNodes: len(graph.Nodes),
			TotalEdges: len(graph.Edges),
		},
		Graph: graph,
	}

	bytes, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		fmt.Printf(`{"error": "Failed to serialize JSON map: %v"}`, err)
		return
	}
	fmt.Println(string(bytes))
}
