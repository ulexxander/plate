package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/ulexxander/plate/commands"
	"github.com/ulexxander/plate/config"
	"github.com/ulexxander/plate/templates/providers"
)

var (
	ErrNoArgs     = errors.New("no arguments provided")
	ErrBadCommand = errors.New("received unknown command")
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.NewDefault()

	provLocal := providers.NewLocal(cfg)
	if err := provLocal.Init(); err != nil {
		return errors.Wrap(err, "error when initializing local templates")
	}

	args := os.Args[1:]
	if len(args) == 0 {
		return ErrNoArgs
	}

	command := args[0]
	commandArgs := args[1:]

	var err error
	switch command {
	case "new":
		c := commands.New{Config: cfg, Provider: provLocal}
		err = c.Exec(commandArgs)
	case "init":
		c := commands.Init{}
		err = c.Exec(commandArgs)
	default:
		err = fmt.Errorf("received unknown command: %s", command)
	}

	return err
}
