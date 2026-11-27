package register

import (
    providers "grovehub/services/integration/internal/app/providers"
    contracts "grovehub/services/integration/internal/app/providers/contracts"
    impl "grovehub/services/integration/internal/app/providers/operations/register/strategies/impl"
)

func NewRegistry(items ...providers.Provider) (contracts.Registry,error) {
    return impl.NewRegistry(items...)
}
