package response

import "grovehub/services/integration/providers"

type ProviderDetails struct {
	Provider         providers.Descriptor `json:"provider"`
	AvailableActions []providers.Action   `json:"availableActions"`
	OAuthConfigured  bool                 `json:"oauthConfigured"`
}
