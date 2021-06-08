package commands

import (
	"errors"
)

type Init struct {
	Args []string
}

func (n *Init) Exec(args []string) error {
	return errors.New("init cmd not implemented")
}
