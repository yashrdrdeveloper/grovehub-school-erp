package search

type Request struct {
	SearchText string      `json:"searchText,omitempty"`
	Filters    []Filter    `json:"filters,omitempty"`
	Sort       []Sort      `json:"sort,omitempty"`
	PageInfo   PageRequest `json:"pageInfo"`
}
