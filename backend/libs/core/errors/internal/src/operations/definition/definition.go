package definition

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors/internal/src/operations/definition/validators"
)

// NewDefinition accepts one optional HTTP error status, defaulting to 500.
func NewDefinition(code contracts.Code, message string, httpStatus ...int) contracts.Definition {
	status := validators.HTTPStatus(httpStatus)

	return contracts.NewDefinition(code, message, status)
}
