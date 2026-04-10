package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lucas/netmap/internal/export"
	"github.com/lucas/netmap/internal/mapper"
	"github.com/lucas/netmap/internal/models"
	"github.com/lucas/netmap/internal/output"
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
	var outFormat string
	var focus string
	var verbose bool

	flag.StringVar(&domain, "d", "", "")
	flag.StringVar(&domain, "domain", "", "")
	flag.StringVar(&mode, "m", "basic", "")
	flag.StringVar(&mode, "mode", "basic", "")
	flag.StringVar(&outFormat, "o", "text", "")
	flag.StringVar(&outFormat, "output", "text", "")
	flag.StringVar(&focus, "f", "all", "")
	flag.StringVar(&focus, "focus", "all", "")
	flag.BoolVar(&verbose, "v", false, "")

	flag.Usage = func() {
		fmt.Printf("%s%s%s\n", output.Cyan, banner, output.Reset)
		fmt.Printf("%sNetMap Intelligence Toolkit - Visual Network Mapper%s\n", output.Bold, output.Reset)
		fmt.Println("\nGebruik:")
		fmt.Println("  netmap -d <target> [flags]")
		fmt.Println("\nFlags:")
		fmt.Println("  -d, --domain string   De website die je wilt analyseren. (Verplicht)")
		fmt.Println("  -f, --focus string    Focus modus: all, auth, admin, api (standaard \"all\")")
		fmt.Println("  -o, --output string    Output formaat: text, json (standaard \"text\")")
		fmt.Println("  -m, --mode string      Mapping modus: basic, advanced (standaard \"basic\")")
		fmt.Println("  -v, --verbose          Toon debug logging (OSINT, HTTP responses, errors)")
		fmt.Println("\nVoorbeelden:")
		fmt.Println("  netmap -d voorbeeld.nl")
		fmt.Println("  netmap -d voorbeeld.nl -f admin")
		fmt.Println("  netmap -d voorbeeld.nl -o json > scan.json\n")
	}

	flag.Parse()

	if domain == "" {
		flag.Usage()
		os.Exit(1)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "[DEBUG] Starting NetMap for target: %s\n", domain)
		fmt.Fprintf(os.Stderr, "[DEBUG] Mode: %s, Output: %s, Focus: %s\n", mode, outFormat, focus)
	}

	// Initialize Mapper Graph
	m := mapper.NewMapper(domain)
	m.Run()

	focusMode := models.FocusMode(focus)

	// Render Output
	if outFormat == "json" || outFormat == "JSON" {
		export.PrintJSON(&m.Graph)
	} else {
		output.PrintTree(&m.Graph, focusMode)
	}
}
