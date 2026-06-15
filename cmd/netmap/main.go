package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/netseries/netmap/internal/export"
	"github.com/netseries/netmap/internal/mapper"
	"github.com/netseries/netmap/internal/models"
	"github.com/netseries/netmap/internal/output"
)

const banner = `
 ███╗   ██╗███████╗████████╗███╗   ███╗ █████╗ ██████╗ 
 ████╗  ██║██╔════╝╚══██╔══╝████╗ ████║██╔══██╗██╔══██╗
 ██╔██╗ ██║█████╗     ██║   ██╔████╔██║███████║██████╔╝
 ██║╚██╗██║██╔══╝     ██║   ██║╚██╔╝██║██╔══██║██╔═══╝ 
 ██║ ╚████║███████╗   ██║   ██║ ╚═╝ ██║██║  ██║██║     
 ╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝     
`

func main() {
	var domain string
	var mode string
	var pack string
	var wordlist string
	var outFormat string
	var focus string
	var verbose bool

	flag.StringVar(&domain, "d", "", "")
	flag.StringVar(&domain, "domain", "", "")
	flag.StringVar(&mode, "m", "basic", "")
	flag.StringVar(&mode, "mode", "basic", "")
	flag.StringVar(&pack, "p", "standard", "")
	flag.StringVar(&pack, "pack", "standard", "")
	flag.StringVar(&wordlist, "w", "", "")
	flag.StringVar(&wordlist, "wordlist", "", "")
	flag.StringVar(&outFormat, "o", "text", "")
	flag.StringVar(&outFormat, "output", "text", "")
	flag.StringVar(&focus, "f", "all", "")
	flag.StringVar(&focus, "focus", "all", "")
	flag.BoolVar(&verbose, "v", false, "")

	flag.Usage = func() {
		fmt.Printf("%s%s%s\n", output.Cyan, banner, output.Reset)
		fmt.Printf("%sNetseries Intelligence Toolkit - Visual Network Mapper%s\n", output.Bold, output.Reset)
		fmt.Println("\nUsage:")
		fmt.Println("  netmap -d <target> [flags]")
		fmt.Println("\nFlags:")
		fmt.Println("  -d, --domain string      The target website/node to analyze. (Required)")
		fmt.Println("  -p, --pack string        Discovery Pack: standard, dns-extended, web-deep, api-focused, admin-stealth, full, ultra")
		fmt.Println("  -w, --wordlist string    Path to a custom discovery wordlist (.txt file)")
		fmt.Println("  -f, --focus string       Focus mode: all, auth, admin, api, config, dev, dns (default \"all\")")
		fmt.Println("  -o, --output string      Output format: text, json (default \"text\")")
		fmt.Println("  -m, --mode string        Mapping mode: basic, advanced (default \"basic\")")
		fmt.Println("  -v, --verbose            Show debug logs (OSINT query, HTTP responses, errors)")
		fmt.Println("\nExamples:")
		fmt.Println("  netmap -d example.com -p ultra")
		fmt.Println("  netmap -d example.com -m advanced -v")
		fmt.Println("  netmap -d example.com -f dns")
		fmt.Println("  netmap -d example.com -o json > map.json")
	}

	flag.Parse()

	if domain == "" {
		flag.Usage()
		os.Exit(1)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "[DEBUG] Starting NetMap for target: %s\n", domain)
		fmt.Fprintf(os.Stderr, "[DEBUG] Mode: %s, Pack: %s, Output: %s, Focus: %s\n", mode, pack, outFormat, focus)
	}

	// Initialize Mapper Graph
	m := mapper.NewMapper(domain, mode, pack, wordlist)
	m.Run()

	focusMode := models.FocusMode(focus)

	// Render Output
	if outFormat == "json" || outFormat == "JSON" {
		export.PrintJSON(&m.Graph)
	} else {
		output.PrintTree(&m.Graph, focusMode)
	}
}
