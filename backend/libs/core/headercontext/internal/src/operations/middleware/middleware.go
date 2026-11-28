package middleware

import (
	"net/http"

	"grovehub/libs/core/headercontext/internal/src/contracts"
	"grovehub/libs/core/headercontext/internal/src/operations/metadata"
	"grovehub/libs/core/headercontext/internal/src/operations/read"
)

// New captures validated metadata before invoking the downstream handler.
func New(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		values, err := read.HTTPHeaders(request.Header)
		if err != nil {
			http.Error(writer, contracts.ErrInvalidHeader.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(writer, request.WithContext(metadata.WithValues(request.Context(), values)))
	})
}
