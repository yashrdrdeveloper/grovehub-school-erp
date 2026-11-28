package operator

type SearchInput struct {
	PageNumber int `json:"pageNumber"`
	RowSize    int `json:"rowSize"`
}

type ReceiptSearchInput struct {
	SearchInput
	SchoolID  string `json:"schoolId,omitempty"`
	Operation string `json:"operation,omitempty"`
}
