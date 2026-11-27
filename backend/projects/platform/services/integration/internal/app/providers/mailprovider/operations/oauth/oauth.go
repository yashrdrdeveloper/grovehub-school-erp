package oauth

import (
    contracts "grovehub/services/integration/internal/app/providers/mailprovider/contracts"
    impl "grovehub/services/integration/internal/app/providers/mailprovider/operations/oauth/strategies/impl"
)

func NewRegistry() contracts.OAuthRegistry { return impl.NewOAuthRegistry() }
