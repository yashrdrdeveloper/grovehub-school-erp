package read

import (
	"net/http"

	"grovehub/libs/core/headercontext/internal/src/constants"
	"grovehub/libs/core/headercontext/internal/src/contracts"
	"grovehub/libs/core/headercontext/internal/src/helpers"
	coreidempotency "grovehub/libs/core/idempotency"
)

// HTTPHeaders reads only supported metadata and validates it before returning any values.
func HTTPHeaders(header http.Header) (contracts.Values, error) {
	idempotency, err := helpers.ReadOne(header, constants.HeaderIdempotencyKey)
	if err != nil {
		return contracts.Values{}, err
	}
	if idempotency != "" {
		idempotency, err = coreidempotency.ValidateKey(idempotency)
		if err != nil {
			return contracts.Values{}, contracts.ErrInvalidHeader
		}
	}
	correlation, err := helpers.ReadOne(header, constants.HeaderCorrelationID)
	if err != nil || !helpers.ValidID(correlation) {
		return contracts.Values{}, contracts.ErrInvalidHeader
	}
	requestID, err := helpers.ReadOne(header, constants.HeaderRequestID)
	if err != nil || !helpers.ValidID(requestID) {
		return contracts.Values{}, contracts.ErrInvalidHeader
	}
	return contracts.Values{IdempotencyKey: idempotency, CorrelationID: correlation, RequestID: requestID}, nil
}
