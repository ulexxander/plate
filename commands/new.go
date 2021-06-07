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
	ErrNoArgsNew    = errors.New("new command requires at least 1 argument")
	ErrInvalidParam = errors.New("invalid tempate param")
)

type New struct {
	Args     []string
	Config   *config.Shape
	Provider *providers.Local // TODO: replace with interface
}

func (n *New) Exec() error {
	if len(n.Args) < 1 {
		return ErrNoArgsNew
	}

	parents := n.Args[:len(n.Args)-1]
	name := n.Args[len(n.Args)-1]

	provided, err := n.Provider.Provide(parents, name)
	if err != nil {
		return err
	}

	tpl, err := template.New(name).Parse(provided)
	if err != nil {
		return errors.Wrap(err, "error when parsing template")
	}

	outf, err := n.openOutputFile(parents, name)
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

func (n *New) openOutputFile(parents []string, name string) (*os.File, error) {
	fname := n.Config.Generated.OutputDir + "/" + name

	outf, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0666)
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

		if line == "end" {
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
