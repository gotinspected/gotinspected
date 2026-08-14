package plugins

// LanguagePlugin defines how a specific programming language is analyzed.
type LanguagePlugin interface {
	Name() string
	Match(filename string) bool
	ShouldExclude(filename string) bool
	IsSignificant(line string) bool
}

// Registry holds all loaded plugins.
type Registry struct {
	plugins       []LanguagePlugin
	defaultPlugin LanguagePlugin
}

// NewRegistry creates a registry with our built-in plugins.
func NewRegistry() *Registry {
	return &Registry{
		plugins: []LanguagePlugin{
			&GoPlugin{},
		},
		defaultPlugin: &DefaultPlugin{},
	}
}

// Get finds the first matching plugin for a file, or returns the default.
func (r *Registry) Get(filename string) LanguagePlugin {
	for _, p := range r.plugins {
		if p.Match(filename) {
			return p
		}
	}
	return r.defaultPlugin
}
