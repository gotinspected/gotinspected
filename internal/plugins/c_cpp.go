package plugins

import (
	"path/filepath"
	"strings"
)

// ==========================================
// C Language Implementation
// ==========================================

type CLanguage struct{}

func (l *CLanguage) Name() string { return "C" }
func (l *CLanguage) Match(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".c" || ext == ".h"
}
func (l *CLanguage) ShouldExclude(filename string) bool { return false }
func (l *CLanguage) NewAnalyzer() FileAnalyzer {
	return &cAnalyzer{}
}

type cAnalyzer struct {
	inBlockComment bool
}

func (a *cAnalyzer) AnalyzeLine(line string) LineType {
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

// ==========================================
// C++ Language Implementation
// ==========================================

type CppLanguage struct{}

func (l *CppLanguage) Name() string { return "C++" }
func (l *CppLanguage) Match(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".cpp" || ext == ".hpp" || ext == ".cc" || ext == ".cxx" || ext == ".hh"
}
func (l *CppLanguage) ShouldExclude(filename string) bool { return false }
func (l *CppLanguage) NewAnalyzer() FileAnalyzer {
	return &cppAnalyzer{}
}

type cppAnalyzer struct {
	inBlockComment bool
}

func (a *cppAnalyzer) AnalyzeLine(line string) LineType {
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
