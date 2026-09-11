package plugins

import "strings"

// --- Factory ---

type DefaultLanguage struct{}

func (l *DefaultLanguage) Name() string                       { return "Default" }
func (l *DefaultLanguage) Match(filename string) bool         { return true }
func (l *DefaultLanguage) ShouldExclude(filename string) bool { return false }
func (l *DefaultLanguage) NewAnalyzer() FileAnalyzer {
	return &defaultAnalyzer{}
}

// --- Stateless Analyzer ---

type defaultAnalyzer struct{}

func (a *defaultAnalyzer) AnalyzeLine(line string) LineType {
	if strings.TrimSpace(line) == "" {
		return TypeEmpty
	}
	return TypeCode
}
