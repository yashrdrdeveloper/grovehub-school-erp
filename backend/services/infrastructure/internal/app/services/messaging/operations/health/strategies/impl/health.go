package impl

import (
    "context"
    healthstrategies "grovehub/services/infrastructure/internal/app/services/messaging/operations/health/strategies"
    providerstrategies "grovehub/services/infrastructure/internal/app/services/messaging/strategies"
)

type healthStrategy struct { provider providerstrategies.Strategy }
var _ healthstrategies.Strategy = (*healthStrategy)(nil)
func New(provider providerstrategies.Strategy) healthstrategies.Strategy { return &healthStrategy{provider:provider} }
func (strategy *healthStrategy) Execute(ctx context.Context) error { return strategy.provider.Health(ctx) }
