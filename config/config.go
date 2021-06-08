package config

type VerboseLevel int

const (
	Silent VerboseLevel = iota
	Error
	Debug
)

type Templates struct {
	LocalDir          string
	ManifestExtension string
}

type Generated struct {
	OutputDir string
}

type Shape struct {
	Templates Templates
	Generated Generated
	Verbosity VerboseLevel
}

func NewDefault() *Shape {
	return &Shape{
		Templates: Templates{
			LocalDir:          "testenv/_templates",
			ManifestExtension: ".plate.json",
		},
		Generated: Generated{
			OutputDir: "testenv/out",
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
