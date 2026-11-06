package services

import (
	"context"

	integrationresponse "grovehub/services/integration/dto/response"
	"grovehub/services/integration/providers"
)

type ProviderService interface {
	Providers(ctx context.Context) []providers.Descriptor
	Provider(ctx context.Context, id string) (integrationresponse.ProviderDetails, error)
}
