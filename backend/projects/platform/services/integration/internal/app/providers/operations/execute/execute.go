package execute

import (
    contracts "grovehub/services/integration/internal/app/providers/contracts"
    impl "grovehub/services/integration/internal/app/providers/operations/execute/strategies/impl"
)

func NewActionRegistry() contracts.ActionRegistry { return impl.NewActionRegistry() }
