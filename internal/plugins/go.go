package plugins

import (
	"path/filepath"
	"strings"
)

// --- Factory ---

type GoLanguage struct{}

func (l *GoLanguage) Name() string               { return "Go" }
func (l *GoLanguage) Match(filename string) bool { return filepath.Ext(filename) == ".go" }
func (l *GoLanguage) ShouldExclude(filename string) bool {
	return strings.HasSuffix(filename, ".pb.go")
}
func (l *GoLanguage) NewAnalyzer() FileAnalyzer {
	return &goAnalyzer{} // Spawn a fresh, stateful analyzer
}

// --- Stateful Analyzer ---

type goAnalyzer struct {
	inBlockComment bool
}

func (a *goAnalyzer) AnalyzeLine(line string) LineType {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return TypeEmpty
	}

	if a.inBlockComment {
		if strings.Contains(trimmed, "*/") {
			a.inBlockComment = false
		}
		return TypeComment
	}

	if strings.HasPrefix(trimmed, "/*") {
		if !strings.Contains(trimmed, "*/") {
			a.inBlockComment = true
		}
		return TypeComment
	}

	if strings.HasPrefix(trimmed, "//") {
		return TypeComment
	}

	return TypeCode
}
