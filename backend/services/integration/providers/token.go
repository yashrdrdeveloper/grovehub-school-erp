package providers

import "time"

// OAuthToken contains only provider-issued authentication material.
// The Integration credential vault encrypts this structure at rest.
type OAuthToken struct {
 AccessToken string `json:"access_token"`
 RefreshToken string `json:"refresh_token"`
 ExpiresIn int `json:"expires_in"`
 TokenType string `json:"token_type"`
 Expiry time.Time `json:"expiry"`
}
