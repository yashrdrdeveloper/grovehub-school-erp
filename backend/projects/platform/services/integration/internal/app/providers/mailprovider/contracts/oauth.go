package contracts

import (
    "context"
    dto "grovehub/services/integration/internal/app/providers/mailprovider/dto"
)

type OAuthStrategy interface {
	AuthorizationURL(state string) string
	Exchange(ctx context.Context, code string) (dto.OAuthToken, error)
	Mailbox(ctx context.Context, token dto.OAuthToken) (string, error)
	Refresh(ctx context.Context, token dto.OAuthToken) (dto.OAuthToken, error)
}

type OAuthRegistry interface {
	Register(provider string, strategy OAuthStrategy) error
	Find(provider string) (OAuthStrategy, error)
}
