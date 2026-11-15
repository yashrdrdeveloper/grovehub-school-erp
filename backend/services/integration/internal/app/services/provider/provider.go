package provider

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    providerdto "grovehub/services/integration/internal/app/dto/provider"
	"context"
)

type ProviderService interface {
	Providers(ctx context.Context) []providerdto.Descriptor
	Provider(ctx context.Context, id string) (providerdto.ProviderDetails, error)
}
