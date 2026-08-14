package plugins

import "strings"

type DefaultPlugin struct{}

func (p *DefaultPlugin) Name() string { return "Default" }
func (p *DefaultPlugin) Match(filename string) bool { return true }
func (p *DefaultPlugin) ShouldExclude(filename string) bool { return false }
func (p *DefaultPlugin) IsSignificant(line string) bool {
	return strings.TrimSpace(line) != ""
}
