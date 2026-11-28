package workload

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/contracts"

	"context"
	"net/http"
	"strings"
	"time"
)

type ScopedWorkloadVerifier struct {
	Credentials contracts.CredentialVerifier
	Audience string
	Routes map[string]contracts.RoutePolicy
}

type workloadContextKey struct{}

func (v ScopedWorkloadVerifier) VerifyWorkload(ctx context.Context, request *http.Request) (context.Context, error) {
	if v.Credentials == nil || v.Audience == "" || request == nil || ctx == nil {
		return nil, contracts.ErrWorkloadDenied
	}
	policy, ok := v.Routes[request.Pattern]
	if !ok || policy.Scope == "" {
		return nil, contracts.ErrWorkloadDenied
	}
	values := request.Header.Values("Authorization")
	if len(values) != 1 {
		return nil, contracts.ErrWorkloadDenied
	}
	scheme, token, ok := strings.Cut(values[0], " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" || strings.ContainsAny(token, " \t\r\n,") {
		return nil, contracts.ErrWorkloadDenied
	}
	claims, err := v.Credentials.VerifyCredential(ctx, token)
	if err != nil || claims.Subject == "" || claims.Audience != v.Audience || !claims.ExpiresAt.After(time.Now()) {
		return nil, contracts.ErrWorkloadDenied
	}
	if (claims.TenantID == "") != (claims.ActorGID == "") || (policy.RequireActor && claims.ActorGID == "") {
		return nil, contracts.ErrWorkloadDenied
	}
	for _, scope := range claims.Scopes {
		if scope == policy.Scope {
			claims.Scopes = append([]string(nil), claims.Scopes...)
			return context.WithValue(ctx, workloadContextKey{}, claims), nil
		}
	}
	return nil, contracts.ErrWorkloadDenied
}

func WorkloadFrom(ctx context.Context) (contracts.WorkloadClaims, bool) {
	if ctx == nil {
		return contracts.WorkloadClaims{}, false
	}
	claims, ok := ctx.Value(workloadContextKey{}).(contracts.WorkloadClaims)
	claims.Scopes = append([]string(nil), claims.Scopes...)
	return claims, ok
}
