package bootstrap

import (
	"log/slog"
	"gorm.io/gorm"
	"grovehub/services/infrastructure/client/messaging"
	platformbootstrap "grovehub/libs/platform/bootstrap"
)

// Dependencies specifies the typed inputs needed to construct the service runtime.
type Dependencies struct {
	Logger *slog.Logger
	DB *gorm.DB
	Events messaging.EventStore
}

type Service = platformbootstrap.RuntimeService[*gorm.DB, Dependencies, Runtime]
