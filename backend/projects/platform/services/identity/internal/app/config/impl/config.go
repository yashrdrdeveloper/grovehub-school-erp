package impl

import (
	"encoding/base64"
	"errors"
	"net/url"
	"net"
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
		AllowedOrigin: os.Getenv("IDENTITY_ALLOWED_ORIGIN"),
	}

	switch os.Getenv("IDENTITY_RBAC_MODE") {
	case "", "strict":
	case "legacy-development":
		settings.LegacyRBACDevelopment = true
	default:
		return config.Config{}, errors.New("invalid IDENTITY_RBAC_MODE")
	}

	// Empty means disabled. A single exact tenant is intentionally required for
	// controlled rollout; there is no wildcard or global execution switch.
	settings.ApprovalExecutionTenant = os.Getenv("IDENTITY_APPROVAL_EXECUTION_TENANT")
	if settings.ApprovalExecutionTenant != "" && (strings.TrimSpace(settings.ApprovalExecutionTenant) != settings.ApprovalExecutionTenant || settings.ApprovalExecutionTenant == "*" || strings.ContainsAny(settings.ApprovalExecutionTenant, ",\x00") || settings.LegacyRBACDevelopment) {
		return config.Config{}, errors.New("invalid IDENTITY_APPROVAL_EXECUTION_TENANT")
	}
	if settings.LegacyRBACDevelopment {
		origin, err := url.Parse(settings.AllowedOrigin)
		if err != nil || origin.Scheme != "http" || !trustedOriginScheme(origin) || os.Getenv("OTEL_DEPLOYMENT_ENVIRONMENT") != "development" {
			return config.Config{}, errors.New("legacy-development requires explicit development environment and HTTP loopback origin")
		}
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
		if err != nil || !trustedOriginScheme(origin) || origin.Host == "" ||
			origin.Path != "" || origin.RawQuery != "" || origin.Fragment != "" ||
			origin.User != nil {
			return config.Config{}, errors.New("identity allowed origin must be HTTPS, or HTTP on loopback")
		}
	}

	return settings, nil
}

// Local Vite development may use HTTP only on a genuine loopback host.
// Public hosts must always use HTTPS.
func trustedOriginScheme(origin *url.URL) bool {
	if origin.Scheme == "https" {
		return true
	}
	if origin.Scheme != "http" {
		return false
	}
	host := origin.Hostname()
	if host == "localhost" {
		return true
	}
	address := net.ParseIP(host)
	return address != nil && address.IsLoopback()
}
