package impl

import (
    "context"
    "grovehub/services/integration/internal/app/providers"
    strategies "grovehub/services/integration/internal/app/services/provider/operations/list/strategies"
)

type listStrategy struct {registry providers.Registry}
var _ strategies.Strategy = (*listStrategy)(nil)
func New(registry providers.Registry) strategies.Strategy {return &listStrategy{registry:registry}}
func (s *listStrategy) Execute(_ context.Context) []providers.Descriptor {return s.registry.List()}
