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
	flag.Parse()

	if domain == "" {
		fmt.Println("Error: Target domain is required.")
		fmt.Println("Usage: netmap -d <domain>")
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
