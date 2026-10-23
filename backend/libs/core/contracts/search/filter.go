package search

type Filter struct {
	Field    string   `json:"field"`
	Operator Operator `json:"operator"`
	Values   []string `json:"values,omitempty"`
}
