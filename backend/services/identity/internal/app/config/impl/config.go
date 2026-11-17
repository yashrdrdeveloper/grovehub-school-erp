package impl

import (
	"encoding/base64"
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"

	"grovehub/services/identity/internal/app/config"
)

func Load() (config.Config, error) {
	settings := config.Config{
		AdminIDs: make(map[string]bool),
		SetupURL: os.Getenv("IDENTITY_SETUP_URL"),
		ResetURL: os.Getenv("IDENTITY_RESET_URL"),
		MailTenantID: os.Getenv("IDENTITY_SYSTEM_MAIL_TENANT_ID"),
		MailConnectionID: os.Getenv("IDENTITY_SYSTEM_MAIL_CONNECTION_ID"),
		AllowedOrigin: os.Getenv("IDENTITY_ALLOWED_ORIGIN"),
	}

	for _, entry := range strings.Split(os.Getenv("IDENTITY_ADMIN_ACCOUNT_IDS"), ",") {
		accountID, err := uuid.Parse(strings.TrimSpace(entry))
		if err == nil {
			settings.AdminIDs[accountID.String()] = true
		}
	}

	key, err := base64.StdEncoding.DecodeString(os.Getenv("IDENTITY_INVITATION_KEY_B64"))
	if err == nil && len(key) == 32 {
		settings.EncryptionKey = key
	}

	if settings.AllowedOrigin != "" {
		origin, err := url.Parse(settings.AllowedOrigin)
		if err != nil || origin.Scheme != "https" || origin.Host == "" ||
			origin.Path != "" || origin.RawQuery != "" || origin.Fragment != "" ||
			origin.User != nil {
			return config.Config{}, errors.New("identity allowed origin must be an HTTPS origin")
		}
	}

	return settings, nil
}
