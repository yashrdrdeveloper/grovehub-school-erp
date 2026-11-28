// Package external owns the internet-facing gateway bridge. Access policy
// is selected from the ServeMux-registered route pattern, not URL substrings.
package impl

import (
	"net/http"
	"regexp"
	"strings"

	coregateway "grovehub/libs/core/gateway"
	"grovehub/launcher/internal/src/services/gateway/operations/external"
	"grovehub/launcher/internal/src/services/gateway/operations/external/helpers"
	"grovehub/launcher/internal/src/services/gateway/operations/ratelimit"
	limiterimpl "grovehub/launcher/internal/src/services/gateway/operations/ratelimit/impl"
	schoolgateway "grovehub/libs/platform/session"
	integrationprincipal "grovehub/services/integration/client/principal"
)

var tierPattern = regexp.MustCompile(`^/[a-z][a-z0-9-]*/([a-z][a-z0-9-]*)(/|$)`)

// @Service
type operation struct {
	mux       *http.ServeMux
	access    schoolgateway.AccessVerifier
	workloads coregateway.WorkloadVerifier
	origin    string
	limiter   ratelimit.Limiter
}

func New(mux *http.ServeMux, access schoolgateway.AccessVerifier, workload coregateway.WorkloadVerifier, allowedOrigin string, limits ...map[string]int) external.Operation {
	configured := map[string]int{}
	if len(limits) > 0 {
		configured = limits[0]
	}

	return &operation{
		mux:       mux,
		access:    access,
		workloads: workload,
		origin:    allowedOrigin,
		limiter:   limiterimpl.New(configured),
	}
}

func (operation *operation) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_, pattern := operation.mux.Handler(request)
	request.Pattern = pattern

	// Only explicitly registered health endpoints may live outside an access tier.
	if pattern == "GET /health" || pattern == "GET /ready" {
		operation.mux.ServeHTTP(writer, request)
		return
	}

	route := pattern
	if strings.HasPrefix(route, request.Method+" ") {
		route = strings.TrimPrefix(route, request.Method+" ")
	}

	matches := tierPattern.FindStringSubmatch(route)
	if len(matches) != 3 {
		helpers.RejectForbidden(writer, request, "route_not_registered_or_unclassified")
		return
	}

	tier := coregateway.Tier(matches[1])
	var serve func(http.ResponseWriter, *http.Request, string)
	switch tier {
	case coregateway.Public:
		serve = operation.servePublic
	case coregateway.Protected:
		serve = operation.serveProtected
	case coregateway.Restricted:
		serve = operation.serveRestricted
	default:
		helpers.RejectForbidden(writer, request, "route_not_registered_or_unclassified")
		return
	}

	if tier != coregateway.Restricted && !helpers.ValidateOrigin(operation.origin, request) {
		helpers.RejectForbidden(writer, request, "origin_mismatch")
		return
	}

	if !operation.limiter.Allow(route, request) {
		writer.Header().Set("Retry-After", "60")
		helpers.Reject(writer, request, http.StatusTooManyRequests)
		return
	}

	serve(writer, request, route)
}

func (operation *operation) servePublic(writer http.ResponseWriter, request *http.Request, _ string) {
	if helpers.RequiresCSRF(request) && !helpers.ValidCSRF(request) {
		helpers.RejectForbidden(writer, request, "csrf_invalid_or_missing")
		return
	}
	operation.mux.ServeHTTP(writer, request)
}

func (operation *operation) serveProtected(writer http.ResponseWriter, request *http.Request, route string) {
	if operation.access == nil {
		helpers.Reject(writer, request, http.StatusUnauthorized)
		return
	}

	token, err := request.Cookie(schoolgateway.AccessCookie)
	if err != nil {
		helpers.Reject(writer, request, http.StatusUnauthorized)
		return
	}

	principal, err := operation.access.VerifyAccess(request.Context(), token.Value)
	if err != nil || principal.AccountID == "" {
		helpers.Reject(writer, request, http.StatusUnauthorized)
		return
	}

	if helpers.RequiresCSRF(request) && !helpers.ValidCSRF(request) {
		helpers.RejectForbidden(writer, request, "csrf_invalid_or_missing")
		return
	}

	ctx := schoolgateway.WithPrincipal(request.Context(), principal)
	if strings.HasPrefix(route, "/identity/protected/") || strings.HasPrefix(route, "/audit/protected/") {
		policy, ok := operation.access.(schoolgateway.ProtectedRouteAuthorizer)
		if !ok || policy.AuthorizeProtectedRoute(ctx, request.Pattern) != nil {
			helpers.RejectForbidden(writer, request, "identity_permission_denied")
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
}

func (operation *operation) serveRestricted(writer http.ResponseWriter, request *http.Request, _ string) {
	// Restricted endpoints must never inherit a browser session.
	if operation.workloads == nil {
		helpers.RejectForbidden(writer, request, "workload_verifier_unavailable")
		return
	}

	ctx, err := operation.workloads.VerifyWorkload(request.Context(), request)
	if err != nil || ctx == nil {
		helpers.RejectForbidden(writer, request, "workload_verification_failed")
		return
	}

	request = request.WithContext(ctx)
	if !helpers.AdaptIdentityProof(request) {
		helpers.RejectForbidden(writer, request, "invalid_identity_proof")
		return
	}

	operation.mux.ServeHTTP(writer, request)
}

var _ external.Operation = (*operation)(nil)
