package helpers

import (
	"log/slog"
	"net/http"

	corehttp "grovehub/libs/core/http"
)

func Reject(writer http.ResponseWriter, request *http.Request, status int) {
	writer.Header().Set("Cache-Control", "no-store")
	// Gateway denies before reaching service controllers. Use the same core
	// writer so the error envelope stays consistent across services.
	_ = corehttp.WriteErrorContext(request.Context(), writer, status, gatewayDefinition(status))
}

// rejectForbidden records a diagnostic reason only in backend logs. The HTTP
// response remains intentionally generic and never reveals policy internals.
func RejectForbidden(writer http.ResponseWriter, request *http.Request, reason string) {
	slog.WarnContext(request.Context(), "gateway request denied",
		"reason", reason,
		"method", request.Method,
		"registered_pattern", request.Pattern,
	)
	Reject(writer, request, http.StatusForbidden)
}
