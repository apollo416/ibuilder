package infra

type Table struct {
	Name        string                    `json:"name"`
	Permissions []*TableFunctionPermision `json:"permissions"`
}

type TableFunctionPermision struct {
	Function   *Function `json:"function"`
	Permisions []string  `json:"permisions"`
}
