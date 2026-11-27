package gateway

import (
	"net/http"
	"regexp"
	"strings"
)

var restrictedPattern = regexp.MustCompile(`^/[a-z][a-z0-9-]*/restricted(/|$)`)

// RestrictedOnly rejects non-restricted routes before invoking the downstream
// gateway. The downstream handler must independently verify workload identity
// and authorize the requested operation. Health endpoints are not exposed here.
func RestrictedOnly(mux *http.ServeMux, downstream http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if mux == nil || downstream == nil {
			denyRestricted(writer)
			return
		}
		_, pattern := mux.Handler(request)
		route := strings.TrimPrefix(pattern, request.Method+" ")
		if !restrictedPattern.MatchString(route) {
			denyRestricted(writer)
			return
		}
		downstream.ServeHTTP(writer, request)
	})
}

func denyRestricted(writer http.ResponseWriter) {
	writer.Header().Set("Cache-Control", "no-store")
	http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}
