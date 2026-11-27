package list

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    "context"
    "grovehub/services/integration/internal/app/services/provider/operations/list/strategies"
)

type ListOperation struct {strategy strategies.Strategy}

func New(strategy strategies.Strategy) *ListOperation {return &ListOperation{strategy:strategy}}

func (service *ListOperation) Providers(ctx context.Context) []providerdto.Descriptor {return service.strategy.Execute(ctx)}
