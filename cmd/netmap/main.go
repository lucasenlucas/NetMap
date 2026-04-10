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
    _   __     __  __  ___          
   / | / /__  / /_/  |/  /___ _____ 
  /  |/ / _ \/ __/ /|_/ / __ ` + "`" + `/ __ \
 / /|  /  __/ /_/ /  / / /_/ / /_/ /
/_/ |_/\___/\__/_/  /_/\__,_/ .___/ 
                           /_/      
`

func main() {
	var domain string
	var mode string
	var outFormat string
	var focus string
	var verbose bool

	flag.StringVar(&domain, "d", "", "Target domain (e.g., example.com)")
	flag.StringVar(&domain, "domain", "", "Target domain (e.g., example.com)")
	flag.StringVar(&mode, "m", "basic", "Mapping mode (basic, advanced)")
	flag.StringVar(&mode, "mode", "basic", "Mapping mode (basic, advanced)")
	flag.StringVar(&outFormat, "o", "text", "Output format (text, json)")
	flag.StringVar(&outFormat, "output", "text", "Output format (text, json)")
	flag.StringVar(&focus, "f", "all", "Focus mode (all, auth, admin, api)")
	flag.StringVar(&focus, "focus", "all", "Focus mode (all, auth, admin, api)")
	flag.BoolVar(&verbose, "v", false, "Verbose output")

	flag.Usage = func() {
		fmt.Printf("%s%s%s\n", output.Cyan, banner, output.Reset)
		fmt.Println("NetMap Intelligence Toolkit - Visual Network Mapper")
		fmt.Println("\nUsage:")
		fmt.Println("  netmap -d <target> [flags]")
		fmt.Println("\nFlags:")
		fmt.Println("  -d, --domain string   Target domain to map (e.g., example.com)")
		fmt.Println("  -f, --focus string    Focus mode: all, auth, admin, api (default \"all\")")
		fmt.Println("  -o, --output string   Output format: text, json (default \"text\")")
		fmt.Println("  -m, --mode string     Mapping mode: basic, advanced (default \"basic\")")
		fmt.Println("  -v, --verbose         Enable verbose debug output")
		fmt.Println("\nExamples:")
		fmt.Println("  netmap -d hackerone.com")
		fmt.Println("  netmap -d example.com -f auth")
		fmt.Println("  netmap -d api.example.com -o json\n")
	}

	flag.Parse()

	if domain == "" {
		flag.Usage()
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("[DEBUG] Starting NetMap for target: %s\n", domain)
		fmt.Printf("[DEBUG] Mode: %s, Output: %s, Focus: %s\n", mode, outFormat, focus)
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
