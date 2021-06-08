package providers

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/ulexxander/plate/config"
	"github.com/ulexxander/plate/templates"
)

var (
	ErrNoSlugs = errors.New("no path to provide template")
)

type Local struct {
	config       *config.Shape
	templatesDir string
	descriptors  map[string]templates.Descriptor
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
		fmt.Println("using templates dir", tplDir)
	}

	l.templatesDir = tplDir
	l.descriptors, err = l.collectDescriptors()
	if err != nil {
		return err
	}

	return nil
}

func (l *Local) collectDescriptors() (map[string]templates.Descriptor, error) {
	descriptors := map[string]templates.Descriptor{}

	err := filepath.Walk(l.templatesDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			// TODO: smarter error handling
			return err
		}

		pathRel := strings.TrimPrefix(path, l.templatesDir+"/")

		if info.IsDir() {
			if l.config.IsDebug() {
				fmt.Println("ignoring directory", pathRel)
			}
			return nil
		}

		if l.isManifestPath(path) {
			if l.config.IsDebug() {
				fmt.Println("ignoring manifest file", pathRel)
			}
			return nil
		}

		descriptor, err := l.makeDescriptor(path, pathRel, info.Name())
		if err != nil {
			return err
		}

		descriptors[descriptor.Slug] = *descriptor
		return nil
	})

	if l.config.IsDebug() {
		fmt.Printf("collected templates %+v\n", descriptors)
	}

	return descriptors, err
}

func (l *Local) isManifestPath(path string) bool {
	return strings.HasSuffix(path, l.config.Templates.ManifestExtension)
}

func (l *Local) makeDescriptor(pathAbs, pathRel, filename string) (*templates.Descriptor, error) {
	ext := filepath.Ext(filename)
	slug := strings.TrimSuffix(pathRel, ext)
	manifestPath := strings.Replace(pathAbs, ext, l.config.Templates.ManifestExtension, 1)

	if _, err := os.Stat(manifestPath); err != nil {
		if os.IsNotExist(err) {
			if l.config.IsDebug() {
				fmt.Println("ignoring template that doesn't have manifest", pathRel)
			}
		} else {
			return nil, errors.Wrap(err, "could not stat file")
		}
	}

	descriptor := templates.Descriptor{
		ManifestPath: manifestPath,
		Path:         pathAbs,
		Slug:         slug,
	}

	return &descriptor, nil
}

func (l *Local) Provide(slug string) (string, error) {
	descriptor, ok := l.descriptors[slug]
	if !ok {
		return "", fmt.Errorf("no templates with slug %s", slug)
	}

	if l.config.IsDebug() {
		fmt.Println("reading template", descriptor.Slug)
	}

	f, err := os.Open(descriptor.Path)
	if err != nil {
		return "", errors.Wrap(err, "could not open template file")
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return "", errors.Wrap(err, "failed to open template file")
	}

	return string(content), nil
}
