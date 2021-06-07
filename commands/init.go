package commands

import "log"

type Init struct {
	Args []string
}

func (n *Init) Exec() error {
	log.Println("init cmd not implemented")
	return nil
}
