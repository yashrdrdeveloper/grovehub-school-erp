package response

// Login is produced only by a configured, durable session issuer.
type Login struct {
    AccessToken string `json:"accessToken"`
    RefreshToken string `json:"refreshToken"`
    TokenType string `json:"tokenType"`
    ExpiresIn int64 `json:"expiresIn"`
}
