package provider

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    "grovehub/services/integration/internal/app/providers"
)
type ProviderDetails struct {
	Provider         providerdto.Descriptor `json:"provider"`
	AvailableActions []providerdto.Action   `json:"availableActions"`
	OAuthConfigured  bool                 `json:"oauthConfigured"`
}
