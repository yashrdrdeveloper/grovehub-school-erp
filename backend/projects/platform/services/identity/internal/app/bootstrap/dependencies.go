package bootstrap

import (
	"log/slog"
	"gorm.io/gorm"
	"grovehub/libs/platform/session"
	person "grovehub/services/foundation/client/person"
	"grovehub/services/infrastructure/client/messaging"
	connection "grovehub/services/integration/client/connection"
	platformbootstrap "grovehub/libs/platform/bootstrap"
)

// Dependencies specifies the typed inputs needed to construct the service runtime.
type Dependencies struct {
	Logger *slog.Logger
	DB *gorm.DB
	Foundation person.Client
	Integration connection.Client
	Events messaging.EventStore
	TenantState []session.TenantStateLookup
}

type Service = platformbootstrap.RuntimeService[*gorm.DB, Dependencies, Runtime]
