package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/yourusername/gitinspector-go/internal/analyzer"
	"github.com/yourusername/gitinspector-go/internal/plugins"
	"github.com/yourusername/gitinspector-go/internal/printer"
)

func main() {
	repoFlag := flag.String("repo", "", "Path to the git repository (optional if passed as argument)")
	flag.Parse()

	// Prioritize the explicit flag; fallback to the first positional argument; default to "."
	rawPath := "."
	if *repoFlag != "" {
		rawPath = *repoFlag
	} else if len(flag.Args()) > 0 {
		rawPath = flag.Args()[0]
	}

	// FIX: Convert whatever the user passed into a clean, absolute path.
	absPath, err := filepath.Abs(rawPath)
	if err != nil {
		log.Fatalf("Failed to resolve path: %v", err)
	}

	// 1. Initialize Plugin Registry
	registry := plugins.NewRegistry()

	// 2. Initialize Analyzer and Printer
	gitAnalyzer := analyzer.New(absPath, registry)
	tablePrinter := printer.NewTablePrinter()

	// 3. Start the spinner (now visually confirming the resolved target path)
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf("  Analyzing Git repository at: %s ...", absPath)
	s.Color("cyan", "bold")
	s.Start()

	// 4. Run Analysis
	result, err := gitAnalyzer.Analyze()
	s.Stop() 

	if err != nil {
		log.Fatalf("\nAnalysis failed for path '%s': %v\nEnsure the target points to or is inside a valid Git repository.", absPath, err)
		os.Exit(1)
	}

	// 5. Render Pretty Table
	tablePrinter.Print(result)
}
