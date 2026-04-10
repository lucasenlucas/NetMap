package output

import (
	"fmt"
	"strings"

	"github.com/lucas/netmap/internal/models"
)

// ANSI color codes for elegant, premium rendering
const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
	Purple = "\033[35m"
	White  = "\033[97m"
)

func PrintTree(graph *models.MapGraph, focus models.FocusMode) {
	var subCount, endpointCount, interestCount, domainCount int

	var roots []*models.Node
	for i, n := range graph.Nodes {
		if n.Type == models.RootNode {
			roots = append(roots, &graph.Nodes[i])
			domainCount++
		} else if n.Type == models.SubdomainNode {
			subCount++
		} else if n.Type == models.EndpointNode {
			endpointCount++
			if isHighInterest(n.Category) {
				interestCount++
			}
		}
	}

	childrenMap := make(map[string][]*models.Node)
	for i, n := range graph.Nodes {
		if n.ParentID != "" {
			childrenMap[n.ParentID] = append(childrenMap[n.ParentID], &graph.Nodes[i])
		}
	}

	// Calculate maximum line length to align tags
	maxLineLen := 0
	for _, root := range roots {
		subs := childrenMap[root.ID]
		for i, sub := range subs {
			isLastSub := i == len(subs)-1
			subPrefix := "│   "
			if isLastSub {
				subPrefix = "    "
			}

			eps := childrenMap[sub.ID]
			for j, ep := range eps {
				if focus == models.FocusAll || string(ep.Category) == string(focus) {
					isLastEp := j == len(eps)-1
					epPrefix := "├──"
					if isLastEp {
						epPrefix = "└──"
					}
					lineStr := fmt.Sprintf("%s%s %s", subPrefix, epPrefix, ep.Label)
					// Handle the exact number of filtered eps effectively by recalibrating isLastEp
					// This is purely for max width estimation, so a loose bound is fine:
					if len(lineStr) > maxLineLen {
						maxLineLen = len(lineStr)
					}
				}
			}
		}
	}

	// Make sure padding exists
	padSpacing := maxLineLen + 4

	// Header
	fmt.Printf("\n%sNetMap%s\n", Bold+Cyan, Reset)
	fmt.Printf("Target: %s\n\n", graph.Target.Domain)

	// Structure
	for _, root := range roots {
		fmt.Printf("%s%s%s\n", Bold, root.Label, Reset)
		subs := childrenMap[root.ID]
		for i, sub := range subs {
			isLastSub := i == len(subs)-1
			prefix := "├──"
			if isLastSub {
				prefix = "└──"
			}
			fmt.Printf("%s %s%s%s\n", prefix, White, sub.Label, Reset)

			eps := childrenMap[sub.ID]
			
			var filteredEps []*models.Node
			for _, ep := range eps {
				if focus == models.FocusAll || string(ep.Category) == string(focus) {
					filteredEps = append(filteredEps, ep)
				}
			}

			subPrefix := "│   "
			if isLastSub {
				subPrefix = "    "
			}

			for j, ep := range filteredEps {
				isLastEp := j == len(filteredEps)-1
				epPrefix := "├──"
				if isLastEp {
					epPrefix = "└──"
				}

				baseStr := fmt.Sprintf("%s%s %s", subPrefix, epPrefix, ep.Label)
				fmt.Print(baseStr)

				// Calculate offset for alignment
				paddingList := padSpacing - len(baseStr)
				if paddingList < 1 {
					paddingList = 1
				}
				fmt.Print(strings.Repeat(" ", paddingList))
				printCategoryLabel(ep.Category)
				fmt.Println()
			}
		}
	}

	// Summary block
	fmt.Printf("\n%sSummary%s\n", Bold, Reset)
	fmt.Printf("  Domains:       %s%d%s\n", White, domainCount, Reset)
	fmt.Printf("  Subdomains:    %s%d%s\n", White, subCount, Reset)
	fmt.Printf("  Endpoints:     %s%d%s\n", White, endpointCount, Reset)
	fmt.Printf("  High-Interest: %s%d%s\n", Yellow, interestCount, Reset)

	if focus != models.FocusAll {
		fmt.Printf("  Filter Focus:  %s%s%s\n", Cyan, focus, Reset)
	}
	fmt.Println()
}

func isHighInterest(cat models.EndpointType) bool {
	return cat == models.TypeAuth || cat == models.TypeAdmin || cat == models.TypeAPI
}

func printCategoryLabel(cat models.EndpointType) {
	switch cat {
	case models.TypeAuth:
		fmt.Printf("%s[AUTH]%s", Purple, Reset)
	case models.TypeAdmin:
		fmt.Printf("%s[ADMIN]%s", Red, Reset)
	case models.TypeAPI:
		fmt.Printf("%s[API]%s", Cyan, Reset)
	case models.TypeGeneral:
		fmt.Printf("%s[GENERAL]%s", Gray, Reset)
	default:
		fmt.Printf("%s[UNKNOWN]%s", Gray, Reset)
	}
}
