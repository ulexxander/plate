package templates

import "errors"

var (
	ErrNoOut = errors.New("out path must be defined for template")
)

type Manifest struct {
	Out    string   `json:"out"`
	Params []string `json:"params"`
}

func (m *Manifest) Validate() error {
	if m.Out == "" {
		return ErrNoOut
	}

	return nil
}
