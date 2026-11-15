package impl

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    providercontracts "grovehub/services/integration/internal/app/providers/contracts"
    "context"
    strategies "grovehub/services/integration/internal/app/services/provider/operations/list/strategies"
)

type listStrategy struct {registry providercontracts.Registry}

var _ strategies.Strategy = (*listStrategy)(nil)

func New(registry providercontracts.Registry) strategies.Strategy {return &listStrategy{registry:registry}}

func (s *listStrategy) Execute(_ context.Context) []providerdto.Descriptor {return s.registry.List()}
