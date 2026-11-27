package provider

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    providerresponse "grovehub/services/integration/internal/app/dto/provider"
	"context"
)

type ProviderService interface {
	Providers(ctx context.Context) []providerdto.Descriptor
	Provider(ctx context.Context, id string) (providerresponse.ProviderDetails, error)
}
