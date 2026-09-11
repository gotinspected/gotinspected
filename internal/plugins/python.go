package plugins

import (
	"path/filepath"
	"strings"
)

type PythonLanguage struct{}

func (l *PythonLanguage) Name() string                       { return "Python" }
func (l *PythonLanguage) Match(filename string) bool         { return filepath.Ext(filename) == ".py" }
func (l *PythonLanguage) ShouldExclude(filename string) bool { return false }
func (l *PythonLanguage) NewAnalyzer() FileAnalyzer          { return &pythonAnalyzer{} }

type pythonAnalyzer struct {
	inDocstring bool
	quoteChar   string
}

func (a *pythonAnalyzer) AnalyzeLine(line string) LineType {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return TypeEmpty
	}

	if a.inDocstring {
		if strings.Contains(trimmed, a.quoteChar) {
			a.inDocstring = false
		}
		return TypeComment
	}

	if strings.HasPrefix(trimmed, `"""`) {
		if len(trimmed) > 3 && strings.HasSuffix(trimmed[3:], `"""`) {
			return TypeComment
		}
		a.inDocstring = true
		a.quoteChar = `"""`
		return TypeComment
	}

	if strings.HasPrefix(trimmed, "'''") {
		if len(trimmed) > 3 && strings.HasSuffix(trimmed[3:], "'''") {
			return TypeComment
		}
		a.inDocstring = true
		a.quoteChar = "'''"
		return TypeComment
	}

	if strings.HasPrefix(trimmed, "#") {
		return TypeComment
	}
	return TypeCode
}
