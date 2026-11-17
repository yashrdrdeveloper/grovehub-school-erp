package oidc

import "context"

// A future OIDC provider has its own protocol-specific interface.
// OIDC login is not password login: never implement unsupported password
// methods with nil return values. Providers must validate PKCE, issuer,
// audience, nonce, state, and provider signature before issuing a session.
type BeginRequest struct {
	RedirectURI string
	Provider string
}

type BeginResult struct {
	AuthorizationURL string
	State string
}

type Claims struct {
	Issuer string
	Subject string
	VerifiedEmail string
}

type Strategy interface {
	Begin(ctx context.Context, request BeginRequest) (BeginResult, error)
	Exchange(ctx context.Context, code string, state string) (Claims, error)
}
