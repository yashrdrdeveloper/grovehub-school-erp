package build

import (
    mailprovider "grovehub/services/integration/internal/app/providers/mailprovider"
    impl "grovehub/services/integration/internal/app/providers/mailprovider/operations/build/strategies/impl"
)

func NewFactory() mailprovider.Factory { return impl.NewFactory() }
