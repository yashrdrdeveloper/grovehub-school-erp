package bootstrap

import (
	"log/slog"
	"gorm.io/gorm"
	"grovehub/libs/platform/session"
	"grovehub/services/infrastructure/client/messaging"
	integrationconfig "grovehub/services/integration/internal/app/config"
	platformbootstrap "grovehub/libs/platform/bootstrap"
)

// Dependencies specifies the typed inputs needed to construct the service runtime.
type Dependencies struct {
	Logger *slog.Logger
	DB *gorm.DB
	Events messaging.MessagingClient
	DevAuth integrationconfig.DevAuth
	TenantState []session.TenantStateLookup
}

type DevAuth = integrationconfig.DevAuth

type Service = platformbootstrap.RuntimeService[*gorm.DB, Dependencies, Runtime]
