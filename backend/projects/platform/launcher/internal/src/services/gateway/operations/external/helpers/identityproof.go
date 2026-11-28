package helpers

import (
	core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	"net/http"
	"strings"
)

// Only the explicitly scoped School workload may supply an identity proof.
// Never interpret browser-supplied identity headers before workload verification.
func AdaptIdentityProof(r *http.Request) bool {
	var cookie, prefix string
	switch r.Pattern {
	case "POST /identity/restricted/auth/refresh":
		cookie, prefix = "__Secure-gh_refresh", ""
	case "GET /identity/restricted/auth/recovery/validate":
		cookie, prefix = "", "Recovery"
	case "GET /identity/restricted/auth/invitations/validate":
		cookie, prefix = "", "Invitation"
	case "POST /identity/restricted/auth/recovery/exchange":
		cookie, prefix = "", "Recovery"
	case "POST /identity/restricted/auth/invitations/exchange":
		cookie, prefix = "", "Invitation"
	case "GET /identity/restricted/auth/reset-password/status":
		cookie, prefix = "__Secure-gh_reset", ""
	case "POST /identity/restricted/auth/reset-password":
		cookie, prefix = "__Secure-gh_reset", ""
	case "GET /identity/restricted/accounts/setup/status":
		cookie, prefix = "__Secure-gh_setup", ""
	case "POST /identity/restricted/accounts/setup":
		cookie, prefix = "__Secure-gh_setup", ""
	default:
		return true
	}

	claims, ok := core.WorkloadFrom(r.Context())
	if !ok || claims.Subject != "school" {
		return false
	}

	values := r.Header.Values("X-GH-Identity-Proof")
	if len(values) != 1 || len(values[0]) == 0 || len(values[0]) > 2048 || strings.ContainsAny(values[0], " \t\r\n,;") {
		return false
	}

	proof := values[0]
	r.Header.Del("X-GH-Identity-Proof")
	r.Header.Del("Cookie")
	r.Header.Del("Authorization")
	if cookie != "" {
		r.AddCookie(&http.Cookie{Name: cookie, Value: proof})
	} else {
		r.Header.Set("Authorization", prefix+" "+proof)
	}

	return true
}
