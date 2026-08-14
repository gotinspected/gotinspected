package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/yourusername/gitinspector-go/internal/analyzer"
	"github.com/yourusername/gitinspector-go/internal/plugins"
	"github.com/yourusername/gitinspector-go/internal/printer"
)

func main() {
	repoDir := flag.String("repo", ".", "Path to the git repository")
	flag.Parse()

	// 1. Initialize Plugin Registry
	registry := plugins.NewRegistry()

	// 2. Initialize Analyzer and Printer
	gitAnalyzer := analyzer.New(*repoDir, registry)
	tablePrinter := printer.NewTablePrinter()

	// 3. Start a Python Rich-style CLI Spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond) // Dot spinner
	s.Suffix = "  Analyzing Git repository trees and applying language filters..."
	s.Color("cyan", "bold")
	s.Start()

	// 4. Run Analysis
	result, err := gitAnalyzer.Analyze()
	s.Stop() // Stop spinner before printing table

	if err != nil {
		log.Fatalf("\nAnalysis failed: %v", err)
		os.Exit(1)
	}

	// 5. Render Pretty Table
	tablePrinter.Print(result)
}
