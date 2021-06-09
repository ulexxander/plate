package commands

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/ulexxander/plate/config"
	"github.com/ulexxander/plate/templates"
	"github.com/ulexxander/plate/templates/providers"
	"github.com/ulexxander/plate/utils/fsutils"
)

var (
	ErrNewNoArgs    = errors.New("new command requires at least 1 argument")
	ErrInvalidParam = errors.New("invalid template param")
)

type New struct {
	Config   *config.Shape
	Provider *providers.Local // TODO: replace with interface
}

func (n *New) Exec(args []string) error {
	if len(args) < 1 {
		return ErrNewNoArgs
	}

	slug := args[0]
	content, err := n.Provider.ProvideContent(slug)
	if err != nil {
		return err
	}

	manifest, err := n.Provider.ProvideManifest(slug)
	if err != nil {
		return err
	}

	tpl, err := template.New(fmt.Sprintf("%s-content", slug)).Parse(content)
	if err != nil {
		return errors.Wrap(err, "error when parsing template")
	}

	input, err := n.scanTemplateInput(manifest)
	if err != nil {
		return err
	}

	outPath, err := n.prepareOutputPath(slug, manifest, input)
	if err != nil {
		return err
	}

	outf, err := n.openOutputFile(outPath)
	if err != nil {
		return err
	}
	defer outf.Close()

	if err := tpl.Execute(outf, input); err != nil {
		return errors.Wrap(err, "failed to execute template")
	}

	return nil
}

func (n *New) scanTemplateInput(manifest *templates.Manifest) (map[string]string, error) {
	params := map[string]string{}

	r := bufio.NewReader(os.Stdin)
	for _, k := range manifest.Params {
		fmt.Printf("%s:\n", k)

		v, err := r.ReadString('\n')
		if err != nil {
			return nil, errors.Wrap(err, "failed to read template param from stdin")
		}

		params[k] = strings.TrimSuffix(v, "\n")
	}

	if n.Config.IsDebug() {
		fmt.Printf("params for template are %+v\n", params)
	}

	return params, nil
}

func (n *New) prepareOutputPath(slug string, manifest *templates.Manifest, input map[string]string) (string, error) {
	t, err := template.New(fmt.Sprintf("%s-out", slug)).Parse(manifest.Out)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse manifest output path")
	}

	var result bytes.Buffer
	if err := t.Execute(&result, input); err != nil {
		return "", errors.Wrap(err, "could not derive output path")
	}

	return result.String(), nil
}

func (n *New) openOutputFile(outPath string) (*os.File, error) {
	outDir := filepath.Dir(outPath)

	if err := fsutils.DirMustExist(outDir); err != nil {
		return nil, errors.Wrap(err, "failed to ensure that output directory exists")
	}

	outf, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return nil, errors.Wrap(err, "could not open output file")
	}

	return outf, nil
}
