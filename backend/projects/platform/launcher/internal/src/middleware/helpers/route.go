package helpers

import (
	"net/http"
	"strings"
)

func Route(request *http.Request, status int) string {
	route := strings.TrimPrefix(request.Pattern, request.Method+" ")
	if route == "" {
		route = "unmatched"
	}

	if route == "/" && request.URL.Path != "/" && status == http.StatusNotFound {
		route = "unmatched"
	}

	return route
}
