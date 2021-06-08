package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/ulexxander/plate/config"
	"github.com/ulexxander/plate/templates/providers"
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
	provided, err := n.Provider.Provide(slug)
	if err != nil {
		return err
	}

	tpl, err := template.New(slug).Parse(provided)
	if err != nil {
		return errors.Wrap(err, "error when parsing template")
	}

	outf, err := n.openOutputFile(slug)
	if err != nil {
		return err
	}
	defer outf.Close()

	input, err := n.scanTemplateInput()
	if err != nil {
		return err
	}

	if err := tpl.Execute(outf, input); err != nil {
		return errors.Wrap(err, "failed to execute template")
	}

	return nil
}

func (n *New) openOutputFile(slug string) (*os.File, error) {
	filename := n.Config.Generated.OutputDir + "/" + slug

	outf, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "error when opening output file")
	}

	return outf, nil
}

func (n *New) scanTemplateInput() (map[string]string, error) {
	params := map[string]string{}

	fmt.Println("enter template parameters key=value")

	inputScanner := bufio.NewScanner(os.Stdin)
	for inputScanner.Scan() {
		line := inputScanner.Text()

		if line == "" {
			break
		}

		// temporary parameters passing
		kv := strings.Split(line, "=")
		if len(kv) < 2 {
			return nil, ErrInvalidParam
		}

		k := kv[0]
		v := kv[1]
		params[k] = v
	}

	if err := inputScanner.Err(); err != nil {
		return nil, errors.Wrap(err, "error when scanning user input")
	}

	return params, nil
}
