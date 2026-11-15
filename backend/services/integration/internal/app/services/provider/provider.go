package provider

import (
    providerdto "grovehub/services/integration/internal/app/dto/provider"
	"context"
	"grovehub/services/integration/internal/app/providers"
)

type ProviderService interface {
	Providers(ctx context.Context) []providers.Descriptor
	Provider(ctx context.Context, id string) (providerdto.ProviderDetails, error)
}
