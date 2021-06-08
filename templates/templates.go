package templates

type Descriptor struct {
	ManifestPath string
	ContentPath  string
	Slug         string
}

type Manifest struct {
	Out    string   `json:"out"`
	Params []string `json:"params"`
}
