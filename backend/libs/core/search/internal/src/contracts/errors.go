package contracts

import "errors"

var (
	ErrInvalidQuery = errors.New("invalid search query")
	ErrInvalidField = errors.New("invalid search field")
	ErrUnsupportedOperator = errors.New("unsupported search operator")
)
