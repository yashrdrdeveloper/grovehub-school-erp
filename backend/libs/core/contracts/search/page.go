package search

const (
	DefaultRows = 25
	MaxRows     = 100
)

type PageRequest struct {
	Rows   int    `json:"rows,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

type PageResponse struct {
	Rows       int    `json:"rows"`
	NextCursor string `json:"nextCursor,omitempty"`
	HasNext    bool   `json:"hasNext"`
}
