package infra

type Api struct {
	Name      string      `json:"name"`
	Resources []*Resource `json:"resources"`
}

type Resource struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Parent  *Resource `json:"parent"`
	Methods []*Method `json:"api_methods"`
}

type Method struct {
	Function  *Function `json:"function"`
	Method    string    `json:"method"`
	Responses []string  `json:"responses"`
}
