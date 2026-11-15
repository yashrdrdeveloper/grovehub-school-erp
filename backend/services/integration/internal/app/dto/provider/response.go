package provider

import "grovehub/services/integration/internal/app/providers"

type ProviderDetails struct {
	Provider         providers.Descriptor `json:"provider"`
	AvailableActions []providers.Action   `json:"availableActions"`
	OAuthConfigured  bool                 `json:"oauthConfigured"`
}
