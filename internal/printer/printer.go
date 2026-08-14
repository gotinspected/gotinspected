package printer

import "github.com/yourusername/gitinspector-go/internal/analyzer"

// Printer defines the contract for outputting analysis results.
type Printer interface {
	Print(result *analyzer.Result)
}
