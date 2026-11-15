package list

import (
    "context"
    "grovehub/services/integration/internal/app/providers"
    "grovehub/services/integration/internal/app/services/provider/operations/list/strategies"
)

type ListOperation struct {strategy strategies.Strategy}

func New(strategy strategies.Strategy) *ListOperation {return &ListOperation{strategy:strategy}}

func (service *ListOperation) Providers(ctx context.Context) []providers.Descriptor {return service.strategy.Execute(ctx)}
