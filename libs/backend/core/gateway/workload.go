package gateway

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

var ErrWorkloadDenied = errors.New("workload authorization denied")

// WorkloadClaims must originate from a cryptographically verified credential,
// never directly from request headers or an unverified token.
type WorkloadClaims struct {
	Subject string
	Audience string
	ExpiresAt time.Time
	Scopes []string
	TenantID string
	ActorGID string
}

// CredentialVerifier verifies credential authenticity, issuer and algorithm.
// ScopedWorkloadVerifier additionally enforces audience, expiry and route scope.
type CredentialVerifier interface {
	VerifyCredential(context.Context, string) (WorkloadClaims, error)
}

type RoutePolicy struct {
	Scope string
	RequireActor bool
}

type ScopedWorkloadVerifier struct {
	Credentials CredentialVerifier
	Audience string
	Routes map[string]RoutePolicy
}

type workloadContextKey struct{}

func (v ScopedWorkloadVerifier) VerifyWorkload(ctx context.Context, request *http.Request) (context.Context, error) {
	if v.Credentials == nil || v.Audience == "" || request == nil || ctx == nil {
		return nil, ErrWorkloadDenied
	}
	policy, ok := v.Routes[request.Pattern]
	if !ok || policy.Scope == "" {
		return nil, ErrWorkloadDenied
	}
	values := request.Header.Values("Authorization")
	if len(values) != 1 {
		return nil, ErrWorkloadDenied
	}
	scheme, token, ok := strings.Cut(values[0], " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" || strings.ContainsAny(token, " \t\r\n,") {
		return nil, ErrWorkloadDenied
	}
	claims, err := v.Credentials.VerifyCredential(ctx, token)
	if err != nil || claims.Subject == "" || claims.Audience != v.Audience || !claims.ExpiresAt.After(time.Now()) {
		return nil, ErrWorkloadDenied
	}
	if (claims.TenantID == "") != (claims.ActorGID == "") || (policy.RequireActor && claims.ActorGID == "") {
		return nil, ErrWorkloadDenied
	}
	for _, scope := range claims.Scopes {
		if scope == policy.Scope {
			claims.Scopes = append([]string(nil), claims.Scopes...)
			return context.WithValue(ctx, workloadContextKey{}, claims), nil
		}
	}
	return nil, ErrWorkloadDenied
}

func WorkloadFrom(ctx context.Context) (WorkloadClaims, bool) {
	if ctx == nil {
		return WorkloadClaims{}, false
	}
	claims, ok := ctx.Value(workloadContextKey{}).(WorkloadClaims)
	claims.Scopes = append([]string(nil), claims.Scopes...)
	return claims, ok
}
