package oauth

import (
    "context"

    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/providers"
    strategies "grovehub/services/integration/internal/app/services/connection/operations/oauth/strategies"
)

type TokenProvider interface {
    ValidProviderToken(ctx context.Context, connectionID, provider string) (providers.OAuthToken,error)
}

type OAuthOperation struct {strategy strategies.Strategy}

func New(strategy strategies.Strategy) *OAuthOperation {return &OAuthOperation{strategy:strategy}}

func (operation *OAuthOperation) BeginOAuth(ctx context.Context, tenantID, id string) (connectiondto.OAuthStart,error) {return operation.strategy.BeginOAuth(ctx,tenantID,id)}

func (operation *OAuthOperation) FinishOAuth(ctx context.Context, provider, code, state string) (connectiondto.OAuthCompletion,error) {return operation.strategy.FinishOAuth(ctx,provider,code,state)}

func (operation *OAuthOperation) ValidProviderToken(ctx context.Context, connectionID, provider string) (providers.OAuthToken,error) {return operation.strategy.ValidProviderToken(ctx,connectionID,provider)}
