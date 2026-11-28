package read

import (
	"net/http"

	httpheaders "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/headers"
	coreidempotency "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext/internal/src/helpers"
)

// HTTPHeaders reads only supported metadata and validates it before returning any values.
func HTTPHeaders(header http.Header) (contracts.Values, error) {
	idempotency, err := helpers.ReadOne(header, httpheaders.IdempotencyKey)
	if err != nil {
		return contracts.Values{}, err
	}
	if idempotency != "" {
		idempotency, err = coreidempotency.ValidateKey(idempotency)
		if err != nil {
			return contracts.Values{}, contracts.ErrInvalidHeader
		}
	}
	correlation, err := helpers.ReadOne(header, httpheaders.CorrelationID)
	if err != nil || !helpers.ValidID(correlation) {
		return contracts.Values{}, contracts.ErrInvalidHeader
	}
	requestID, err := helpers.ReadOne(header, httpheaders.RequestID)
	if err != nil || !helpers.ValidID(requestID) {
		return contracts.Values{}, contracts.ErrInvalidHeader
	}
	return contracts.Values{IdempotencyKey: idempotency, CorrelationID: correlation, RequestID: requestID}, nil
}
