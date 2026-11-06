package providers

import "context"

type OAuthStrategy interface {
	AuthorizationURL(state string) string
	Exchange(ctx context.Context, code string) (OAuthToken, error)
	Mailbox(ctx context.Context, token OAuthToken) (string, error)
	Refresh(ctx context.Context, token OAuthToken) (OAuthToken, error)
}
type OAuthRegistry interface {
	Register(provider string, strategy OAuthStrategy) error
	Find(provider string) (OAuthStrategy, error)
}
