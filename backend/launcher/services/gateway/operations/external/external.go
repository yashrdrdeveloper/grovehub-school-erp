// Package external owns the internet-facing gateway bridge. Access policy
// is selected from the ServeMux-registered route pattern, not URL substrings.
package external

import (
	"crypto/subtle"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	coregateway "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	httpresponse "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/response"
	schoolgateway "grovehub/libs/schoolgateway"
	integrationprincipal "grovehub/services/integration/client/principal"
)

var tierPattern = regexp.MustCompile(`^/[a-z][a-z0-9-]*/(public|protected|restricted)(/|$)`)

type Operation struct {
	mux       *http.ServeMux
	access    schoolgateway.AccessVerifier
	workloads coregateway.WorkloadVerifier
	origin    string
	limiter   *Limiter
}

func New(mux *http.ServeMux, access schoolgateway.AccessVerifier, workload coregateway.WorkloadVerifier, allowedOrigin string, limits ...map[string]int) *Operation {
	configured := map[string]int{}
	if len(limits) > 0 {
		configured = limits[0]
	}
	return &Operation{
		mux:       mux,
		access:    access,
		workloads: workload,
		origin:    allowedOrigin,
		limiter:   NewLimiter(configured),
	}
}

func (operation *Operation) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_, pattern := operation.mux.Handler(request)
	request.Pattern = pattern

	// Only explicitly registered health endpoints may live outside an access tier.
	if pattern == "GET /health" {
		operation.mux.ServeHTTP(writer, request)
		return
	}

	route := pattern
	if strings.HasPrefix(route, request.Method+" ") {
		route = strings.TrimPrefix(route, request.Method+" ")
	}
	matches := tierPattern.FindStringSubmatch(route)
	if len(matches) != 3 {
		operation.rejectForbidden(writer, request, "route_not_registered_or_unclassified")
		return
	}

	if coregateway.Tier(matches[1]) != coregateway.Restricted && !operation.validateOrigin(request) {
		operation.rejectForbidden(writer, request, "origin_mismatch")
		return
	}

	if !operation.limiter.Allow(route, request) {
		writer.Header().Set("Retry-After", "60")
		operation.reject(writer, request, http.StatusTooManyRequests)
		return
	}

	switch coregateway.Tier(matches[1]) {
	case coregateway.Public:
		if requiresCSRF(request) && !validCSRF(request) {
			operation.rejectForbidden(writer, request, "csrf_invalid_or_missing")
			return
		}
		operation.mux.ServeHTTP(writer, request)

	case coregateway.Protected:
		if operation.access == nil {
			operation.reject(writer, request, http.StatusUnauthorized)
			return
		}
		token, err := request.Cookie(schoolgateway.AccessCookie)
		if err != nil {
			operation.reject(writer, request, http.StatusUnauthorized)
			return
		}
		principal, err := operation.access.VerifyAccess(request.Context(), token.Value)
		if err != nil || principal.AccountID == "" {
			operation.reject(writer, request, http.StatusUnauthorized)
			return
		}
		if requiresCSRF(request) && !validCSRF(request) {
			operation.rejectForbidden(writer, request, "csrf_invalid_or_missing")
			return
		}

		ctx := schoolgateway.WithPrincipal(request.Context(), principal)
		if strings.HasPrefix(route, "/identity/protected/") {
			policy, ok := operation.access.(schoolgateway.ProtectedRouteAuthorizer)
			if !ok || policy.AuthorizeProtectedRoute(ctx, pattern) != nil {
				operation.rejectForbidden(writer, request, "identity_permission_denied")
				return
			}
		}
		if principal.TenantID != "" && principal.PersonGID != "" {
			ctx = integrationprincipal.WithVerifiedPrincipal(ctx, integrationprincipal.VerifiedPrincipal{
				TenantID: principal.TenantID,
				ActorGID: principal.PersonGID,
			})
		}
		operation.mux.ServeHTTP(writer, request.WithContext(ctx))

	case coregateway.Restricted:
		// Restricted endpoints must never inherit a browser session.
		if operation.workloads == nil {
			operation.rejectForbidden(writer, request, "workload_verifier_unavailable")
			return
		}
		ctx, err := operation.workloads.VerifyWorkload(request.Context(), request)
		if err != nil || ctx == nil {
			operation.rejectForbidden(writer, request, "workload_verification_failed")
			return
		}
		operation.mux.ServeHTTP(writer, request.WithContext(ctx))

	default:
		operation.rejectForbidden(writer, request, "unsupported_access_tier")
	}
}

func (operation *Operation) validateOrigin(request *http.Request) bool {
	if request.Method == http.MethodGet || request.Method == http.MethodHead {
		return true
	}
	// An empty Origin is not treated as proof of a same-origin browser request.
	if operation.origin != "" && request.Header.Get("Origin") == operation.origin {
		return true
	}
	// The local Vite proxy can be opened as either localhost or 127.0.0.1.
	// Treat only these two loopback spellings as equivalent, on the same port
	// and scheme, when the configured allowed origin itself is loopback.
	allowed, allowedErr := url.Parse(operation.origin)
	incoming, incomingErr := url.Parse(request.Header.Get("Origin"))
	if allowedErr != nil || incomingErr != nil || allowed.Scheme != "http" || incoming.Scheme != "http" || allowed.Port() == "" || allowed.Port() != incoming.Port() {
		return false
	}
	isLoopback := func(host string) bool { return host == "localhost" || host == "127.0.0.1" }
	return isLoopback(allowed.Hostname()) && isLoopback(incoming.Hostname())
}

func requiresCSRF(request *http.Request) bool {
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

func validCSRF(request *http.Request) bool {
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

func (operation *Operation) reject(writer http.ResponseWriter, request *http.Request, status int) {
	writer.Header().Set("Cache-Control", "no-store")
	// Gateway denies before reaching service controllers. Use the same core
	// writer so the error envelope stays consistent across services.
	_ = httpresponse.WriteErrorContext(request.Context(), writer, status, gatewayDefinition(status))
}

// rejectForbidden records a diagnostic reason only in backend logs. The HTTP
// response remains intentionally generic and never reveals policy internals.
func (operation *Operation) rejectForbidden(writer http.ResponseWriter, request *http.Request, reason string) {
	slog.WarnContext(request.Context(), "gateway request denied",
		"reason", reason,
		"method", request.Method,
		"registered_pattern", request.Pattern,
	)
	operation.reject(writer, request, http.StatusForbidden)
}
