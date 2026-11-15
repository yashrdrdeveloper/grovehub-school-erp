// Package impl composes Provider discovery operations.
package impl

import (
    providercontracts "grovehub/services/integration/internal/app/providers/contracts"
    mailcontracts "grovehub/services/integration/internal/app/providers/mailprovider/contracts"
    services "grovehub/services/integration/internal/app/services/provider"
    get "grovehub/services/integration/internal/app/services/provider/operations/get"
    getimpl "grovehub/services/integration/internal/app/services/provider/operations/get/strategies/impl"
    list "grovehub/services/integration/internal/app/services/provider/operations/list"
    listimpl "grovehub/services/integration/internal/app/services/provider/operations/list/strategies/impl"
)

// @Service
type ProviderService struct {
    *list.ListOperation
    *get.GetOperation
}

var _ services.ProviderService = (*ProviderService)(nil)

func NewProviderService(registry providercontracts.Registry, actions providercontracts.ActionRegistry, oauth mailcontracts.OAuthRegistry) services.ProviderService {
    return &ProviderService{
        ListOperation:list.New(listimpl.New(registry)),
        GetOperation:get.New(getimpl.New(registry,actions,oauth)),
    }
}
