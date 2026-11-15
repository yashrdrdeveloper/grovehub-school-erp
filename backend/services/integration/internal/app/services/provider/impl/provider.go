// Package impl composes Provider discovery operations.
package impl

import (
    "grovehub/services/integration/internal/app/providers"
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

func NewProviderService(registry providers.Registry, actions providers.ActionRegistry, oauth providers.OAuthRegistry) services.ProviderService {
    return &ProviderService{
        ListOperation:list.New(listimpl.New(registry)),
        GetOperation:get.New(getimpl.New(registry,actions,oauth)),
    }
}
