package plugins

import (
	"path/filepath"
	"strings"
)

type GoPlugin struct{}

func (p *GoPlugin) Name() string { return "Go" }

func (p *GoPlugin) Match(filename string) bool {
	return filepath.Ext(filename) == ".go"
}

func (p *GoPlugin) ShouldExclude(filename string) bool {
	// Example: Exclude generated protobufs or mock files
	if strings.HasSuffix(filename, ".pb.go") {
		return true
	}
	return false
}

func (p *GoPlugin) IsSignificant(line string) bool {
	line = strings.TrimSpace(line)
	if line == "" {
		return false
	}
	if strings.HasPrefix(line, "//") {
		return false
	}
	return true
}
