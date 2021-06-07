package config

type VerboseLevel int

const (
	Silent VerboseLevel = iota
	Error
	Debug
)

type Shape struct {
	Templates struct {
		LocalDir string
	}
	Generated struct {
		OutputDir string
	}
	Verbosity VerboseLevel
}

func NewDefault() *Shape {
	return &Shape{
		Templates: struct{ LocalDir string }{
			LocalDir: "testenv/_templates",
		},
		Generated: struct{ OutputDir string }{
			OutputDir: ".",
		},
		Verbosity: Debug,
	}
}

func (s *Shape) IsError() bool {
	return s.Verbosity >= Error
}

func (s *Shape) IsDebug() bool {
	return s.Verbosity >= Debug
}
