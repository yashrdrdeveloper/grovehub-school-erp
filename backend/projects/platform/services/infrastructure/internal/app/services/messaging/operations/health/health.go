package health

import (
    "context"
    "grovehub/services/infrastructure/internal/app/services/messaging/operations/health/strategies"
)

type HealthOperation struct { strategy strategies.Strategy }
func New(strategy strategies.Strategy) *HealthOperation { return &HealthOperation{strategy:strategy} }
func (operation *HealthOperation) Health(ctx context.Context) error { return operation.strategy.Execute(ctx) }
