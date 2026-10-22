package response

type Success[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

func NewSuccess[T any](data T) Success[T] {
	return Success[T]{
		Success: true,
		Data:    data,
	}
}
