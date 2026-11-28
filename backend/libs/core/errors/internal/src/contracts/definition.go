package contracts

import "net/http"

// Definition describes a reusable application error without a cause or metadata.
type Definition struct {
	code       Code
	message    string
	httpStatus int
}

// NewDefinition constructs a definition from a status validated by its operation.
func NewDefinition(code Code, message string, status int) Definition {
	return Definition{code: code, message: message, httpStatus: status}
}

func (definition Definition) Code() Code {
	return definition.code
}

func (definition Definition) Message() string {
	return definition.message
}

// HTTPStatus defaults zero-valued definitions to HTTP 500.
func (definition Definition) HTTPStatus() int {
	if definition.httpStatus == 0 {
		return http.StatusInternalServerError
	}

	return definition.httpStatus
}
