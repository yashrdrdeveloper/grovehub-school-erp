package config

type Config struct {
	AdminIDs       map[string]bool
	LegacyRBACDevelopment bool
	ApprovalExecutionTenant string
	EncryptionKey  []byte
	SetupURL       string
	ResetURL       string
	AllowedOrigin  string
}
