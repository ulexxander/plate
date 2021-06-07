package providers

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/ulexxander/plate/config"
)

var (
	ErrNoSlugs = errors.New("no path to provide template")
)

type Local struct {
	config       *config.Shape
	templatesDir string
}

func NewLocal(c *config.Shape) *Local {
	return &Local{
		config: c,
	}
}

func (l *Local) Init() error {
	cwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "could not get current working directory")
	}

	tplDir := filepath.Join(cwd, l.config.Templates.LocalDir)

	if l.config.IsDebug() {
		log.Println("using templates dir", tplDir)
	}

	l.templatesDir = tplDir

	return nil
}

func (l *Local) Provide(parents []string, name string) (string, error) {
	pathRelative := strings.Join(parents, "/") + name
	fname := l.templatesDir + "/" + pathRelative
	if l.config.IsDebug() {
		log.Println("reading template", pathRelative)
	}

	f, err := os.Open(fname)
	if err != nil {
		return "", errors.Wrap(err, "could not open template file")
	}
	defer f.Close()

	fcontent, err := io.ReadAll(f)
	if err != nil {
		return "", errors.Wrap(err, "failed to open template file")
	}

	return string(fcontent), nil
}
