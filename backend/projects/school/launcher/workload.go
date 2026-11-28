package launcher

import (
	"net/http"
	"strings"

	auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/identity/client/authentication"
	business "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/platform/client"
	core "grovehub/libs/core/gateway"
)

type workloadPolicy func(*http.Request) (core.WorkloadClaims, error)

func scoped(scope string) (core.WorkloadClaims, error) {
	return core.WorkloadClaims{Scopes: []string{scope}}, nil
}

func validPolicyRequest(r *http.Request) bool {
	return r != nil && r.URL != nil && r.URL.Fragment == "" && r.URL.RawPath == ""
}

func identityPolicy(r *http.Request) (core.WorkloadClaims, error) {
	if !validPolicyRequest(r) || r.URL.RawQuery != "" {
		return core.WorkloadClaims{}, core.ErrWorkloadDenied
	}
	switch r.Method + " " + r.URL.Path {
	case "POST /identity/restricted/auth/login":
		return scoped("identity.auth.login")
	case "POST /identity/restricted/auth/forgot-password":
		return scoped("identity.auth.forgot")
	case "POST /identity/restricted/auth/refresh":
		return scoped("identity.auth.refresh")
	case "GET /identity/restricted/auth/recovery/validate":
		return scoped("identity.auth.recovery.validate")
	case "GET /identity/restricted/auth/invitations/validate":
		return scoped("identity.auth.invitation.validate")
	case "POST /identity/restricted/auth/recovery/exchange":
		return scoped("identity.auth.recovery.exchange")
	case "POST /identity/restricted/auth/invitations/exchange":
		return scoped("identity.auth.invitation.exchange")
	case "GET /identity/restricted/auth/reset-password/status":
		return scoped("identity.auth.reset.status")
	case "POST /identity/restricted/auth/reset-password":
		return scoped("identity.auth.reset")
	case "GET /identity/restricted/accounts/setup/status":
		return scoped("identity.auth.setup.status")
	case "POST /identity/restricted/accounts/setup":
		return scoped("identity.auth.setup")
	case "POST /identity/restricted/auth/logout":
		return scoped("identity.auth.logout")
	case "POST /identity/restricted/auth/logout-all":
		return scoped("identity.auth.logout.all")
	case "POST /identity/restricted/auth/change-password":
		return scoped("identity.auth.change.password")
	case "POST /identity/restricted/auth/reauthenticate":
		return scoped("identity.auth.reauthenticate")
	case "GET /identity/restricted/auth/sessions":
		return scoped("identity.auth.sessions")
	case "GET /identity/restricted/auth/me":
		return scoped("identity.auth.me")
	}
	if r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/identity/restricted/auth/sessions/") && auth.ValidSessionID(strings.TrimPrefix(r.URL.Path, "/identity/restricted/auth/sessions/")) {
		return scoped("identity.auth.sessions.revoke")
	}
	return businessPolicy(r, "identity.")
}

func auditPolicy(r *http.Request) (core.WorkloadClaims, error) {
	return businessPolicy(r, "audit.")
}

func integrationPolicy(r *http.Request) (core.WorkloadClaims, error) {
	if !validPolicyRequest(r) {
		return core.WorkloadClaims{}, core.ErrWorkloadDenied
	}
	switch r.Method + " " + r.URL.Path {
	case "GET /integrations/restricted/connections/oauth/gmail/callback", "GET /integrations/restricted/connections/oauth/outlook/callback":
		q := r.URL.Query()
		if len(q) == 2 && len(q["code"]) == 1 && len(q["state"]) == 1 && q.Get("code") != "" && q.Get("state") != "" {
			return scoped("integration.oauth.callback")
		}
		return core.WorkloadClaims{}, core.ErrWorkloadDenied
	}
	return businessPolicy(r, "integration.")
}

func businessPolicy(r *http.Request, prefix string) (core.WorkloadClaims, error) {
	if !validPolicyRequest(r) || r.URL.RawQuery != "" {
		return core.WorkloadClaims{}, core.ErrWorkloadDenied
	}
	if scope, ok := business.Scope(r); ok && strings.HasPrefix(scope, prefix) {
		return scoped(scope)
	}
	return core.WorkloadClaims{}, core.ErrWorkloadDenied
}
