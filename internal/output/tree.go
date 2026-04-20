package output

import (
	"fmt"
	"strings"

	"github.com/netseries/netmap/internal/models"
)

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
	Purple = "\033[35m"
	Green  = "\033[32m"
	White  = "\033[97m"
	Blue   = "\033[34m"
)

func PrintTree(graph *models.MapGraph, focus models.FocusMode) {
	var subCount, endpointCount, interestCount, domainCount, dnsCount int

	var roots []*models.Node
	for i, n := range graph.Nodes {
		if n.Type == models.RootNode {
			roots = append(roots, &graph.Nodes[i])
			domainCount++
		} else if n.Type == models.SubdomainNode {
			subCount++
		} else if n.Type == models.EndpointNode {
			if n.Category == models.TypeDNS {
				dnsCount++
			} else {
				endpointCount++
			}
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

	// Dynamic padding
	maxLineLen := 40 
	padSpacing := maxLineLen + 10

	fmt.Printf("\n%sNetMap Intelligence%s\n", Bold+Cyan, Reset)
	fmt.Printf("Target: %s\n\n", graph.Target.Domain)

	for _, root := range roots {
		fmt.Printf("%s%s%s\n", Bold, root.Label, Reset)
		
		var rootSubs []*models.Node
		var rootEps []*models.Node
		for _, child := range childrenMap[root.ID] {
			if child.Type == models.SubdomainNode {
				rootSubs = append(rootSubs, child)
			} else if child.Type == models.EndpointNode {
				if shouldShow(child, focus) {
					rootEps = append(rootEps, child)
				}
			}
		}

		for i, sub := range rootSubs {
			isLastSub := (i == len(rootSubs)-1) && len(rootEps) == 0
			prefix := "├──"
			if isLastSub {
				prefix = "└──"
			}
			fmt.Printf("%s %s%s%s\n", prefix, White, sub.Label, Reset)

			eps := childrenMap[sub.ID]
			var filteredEps []*models.Node
			for _, ep := range eps {
				if shouldShow(ep, focus) {
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

				label := ep.Label
				if len(label) > 60 { label = label[:57] + "..." }
				baseStr := fmt.Sprintf("%s%s %s", subPrefix, epPrefix, label)
				fmt.Print(baseStr)

				paddingList := padSpacing - len(baseStr)
				if paddingList < 1 { paddingList = 1 }
				fmt.Print(strings.Repeat(" ", paddingList))
				printCategoryLabel(ep.Category, ep.Label)
				fmt.Println()
			}
		}

		for j, ep := range rootEps {
			isLastEp := j == len(rootEps)-1
			epPrefix := "├──"
			if isLastEp {
				epPrefix = "└──"
			}
			label := ep.Label
			if len(label) > 60 { label = label[:57] + "..." }
			baseStr := fmt.Sprintf("%s %s", epPrefix, label)
			fmt.Print(baseStr)

			paddingList := padSpacing - len(baseStr)
			if paddingList < 1 { paddingList = 1 }
			fmt.Print(strings.Repeat(" ", paddingList))
			printCategoryLabel(ep.Category, ep.Label)
			fmt.Println()
		}
	}

	fmt.Printf("\n%sSummary%s\n", Bold, Reset)
	fmt.Printf("  Domains:       %s%d%s\n", White, domainCount, Reset)
	fmt.Printf("  Subdomains:    %s%d%s\n", White, subCount, Reset)
	fmt.Printf("  Endpoints:     %s%d%s\n", White, endpointCount, Reset)
	fmt.Printf("  DNS Records:   %s%d%s\n", Blue, dnsCount, Reset)
	fmt.Printf("  High-Interest: %s%d%s\n", Yellow, interestCount, Reset)
	if focus != models.FocusAll {
		fmt.Printf("  Filter Focus:  %s%s%s\n", Cyan, focus, Reset)
	}
	fmt.Println()
}

func shouldShow(node *models.Node, focus models.FocusMode) bool {
	if focus == models.FocusAll {
		return true
	}
	return string(node.Category) == string(focus)
}

func isHighInterest(cat models.EndpointType) bool {
	return cat == models.TypeAuth || cat == models.TypeAdmin || cat == models.TypeAPI || cat == models.TypeConfig || cat == models.TypeDev || cat == models.TypeDNS
}

func printCategoryLabel(cat models.EndpointType, label string) {
	switch cat {
	case models.TypeAuth:
		fmt.Printf("%s[AUTH]%s", Purple, Reset)
	case models.TypeAdmin:
		fmt.Printf("%s[ADMIN]%s", Red, Reset)
	case models.TypeAPI:
		fmt.Printf("%s[API]%s", Cyan, Reset)
	case models.TypeConfig:
		fmt.Printf("%s[CONFIG]%s", Yellow, Reset)
	case models.TypeDev:
		fmt.Printf("%s[DEV]%s", Green, Reset)
	case models.TypeDNS:
		if strings.HasPrefix(label, "MX:") {
			fmt.Printf("%s[MX]%s", Blue, Reset)
		} else if strings.HasPrefix(label, "CNAME:") {
			fmt.Printf("%s[CNAME]%s", Cyan, Reset)
		} else {
			fmt.Printf("%s[TXT]%s", Yellow, Reset)
		}
	case models.TypeGeneral:
		fmt.Printf("%s[GENERAL]%s", Gray, Reset)
	default:
		fmt.Printf("%s[UNKNOWN]%s", Gray, Reset)
	}
}
