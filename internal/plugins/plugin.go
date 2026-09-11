package plugins

type LineType int

const (
	TypeCode LineType = iota
	TypeEmpty
	TypeComment
)

// FileAnalyzer encapsulates the state for parsing a single stream of lines.
type FileAnalyzer interface {
	AnalyzeLine(line string) LineType
}

// Language defines metadata and acts as a factory for its FileAnalyzer.
type Language interface {
	Name() string
	Match(filename string) bool
	ShouldExclude(filename string) bool
	NewAnalyzer() FileAnalyzer
}

type Registry struct {
	languages       []Language
	defaultLanguage Language
}

func NewRegistry() *Registry {
	return &Registry{
		languages: []Language{
			&GoLanguage{},
			&PythonLanguage{},
			&CLanguage{},
			&CppLanguage{},
		},
		defaultLanguage: &DefaultLanguage{},
	}
}

func (r *Registry) Get(filename string) Language {
	for _, l := range r.languages {
		if l.Match(filename) {
			return l
		}
	}
	return r.defaultLanguage
}
