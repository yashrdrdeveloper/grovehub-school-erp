package helpers

import (
	"crypto/subtle"
	schoolgateway "grovehub/internal/session"
	"net/http"
	"strings"
)

func RequiresCSRF(request *http.Request) bool {
	if request.Method == http.MethodGet || request.Method == http.MethodHead {
		return false
	}

	route := strings.TrimPrefix(request.Pattern, request.Method+" ")
	if strings.Contains(route, "/protected/") {
		return true
	}

	switch route {
	case "/identity/public/auth/refresh",
		"/identity/public/auth/reset-password",
		"/identity/public/accounts/setup":
		return true
	default:
		return false
	}
}

func ValidCSRF(request *http.Request) bool {
	cookie, err := request.Cookie(schoolgateway.CSRFName)
	if err != nil {
		return false
	}

	value := request.Header.Get("X-CSRF-Token")
	if len(value) == 0 || len(value) > 256 || len(cookie.Value) == 0 {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(value)) == 1
}
