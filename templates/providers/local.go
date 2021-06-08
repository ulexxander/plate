package providers

import (
	"encoding/json"
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

	// TODO: better templates scanning, go through manifests, not actual templtes
	err := filepath.Walk(l.templatesDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			// TODO: smarter error handling
			return err
		}

		// TODO: do we have safer way? (may fail on shindows)
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
		ContentPath:  pathAbs,
		Slug:         slug,
	}

	return &descriptor, nil
}

func (l *Local) ProvideContent(slug string) (string, error) {
	descriptor, ok := l.descriptors[slug]
	if !ok {
		return "", fmt.Errorf("no template with slug %s", slug)
	}

	if l.config.IsDebug() {
		fmt.Println("reading content of", descriptor.Slug)
	}

	f, err := os.Open(descriptor.ContentPath)
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

func (l *Local) ProvideManifest(slug string) (*templates.Manifest, error) {
	descriptor, ok := l.descriptors[slug]
	if !ok {
		return nil, fmt.Errorf("no template with slug %s", slug)
	}

	if l.config.IsDebug() {
		fmt.Println("reading manifest for", descriptor.Slug)
	}

	f, err := os.Open(descriptor.ManifestPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not open manifest file")
	}
	defer f.Close()

	var manifest templates.Manifest
	if err := json.NewDecoder(f).Decode(&manifest); err != nil {
		return nil, errors.Wrap(err, "failed to decode manifest as json")
	}

	if l.config.IsDebug() {
		fmt.Printf("got manifest for %s %+v\n", descriptor.Slug, manifest)
	}

	if err := manifest.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid manifest")
	}

	return &manifest, nil
}
