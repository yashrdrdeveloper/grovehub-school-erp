package config

type Config struct {
	AdminIDs       map[string]bool
	EncryptionKey  []byte
	SetupURL       string
	ResetURL       string
	AllowedOrigin  string
}
