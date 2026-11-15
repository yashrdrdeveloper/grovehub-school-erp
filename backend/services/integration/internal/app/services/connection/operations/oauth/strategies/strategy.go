package strategies

import (
    "context"

    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/providers"
)

// @Strategy
type Strategy interface {
    BeginOAuth(ctx context.Context, tenantID, id string) (connectiondto.OAuthStart,error)
    FinishOAuth(ctx context.Context, provider, code, state string) (connectiondto.OAuthCompletion,error)
    ValidProviderToken(ctx context.Context, connectionID, provider string) (providers.OAuthToken,error)
}
