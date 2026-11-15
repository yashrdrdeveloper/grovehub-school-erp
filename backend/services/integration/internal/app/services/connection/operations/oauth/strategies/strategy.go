package strategies

import (
    maildto "grovehub/services/integration/internal/app/providers/mailprovider/dto"
    "context"

    connectiondto "grovehub/services/integration/internal/app/dto/connection"
)

// @Strategy
type Strategy interface {
    BeginOAuth(ctx context.Context, tenantID, id string) (connectiondto.OAuthStart,error)
    FinishOAuth(ctx context.Context, provider, code, state string) (connectiondto.OAuthCompletion,error)
    ValidProviderToken(ctx context.Context, connectionID, provider string) (maildto.OAuthToken,error)
}
