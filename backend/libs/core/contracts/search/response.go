package search

type Result[T any] struct {
	Items    []T          `json:"items"`
	PageInfo PageResponse `json:"pageInfo"`
}
