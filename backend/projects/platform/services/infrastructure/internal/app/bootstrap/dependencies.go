package bootstrap

import (
	"log/slog"
	"gorm.io/gorm"
	"grovehub/services/infrastructure/client/messaging"
	platformbootstrap "grovehub/libs/platform/bootstrap"
)

// Dependencies specifies the typed inputs needed to construct the service runtime.
type Dependencies struct {
	Config messaging.Config
	DB *gorm.DB
	Logger *slog.Logger
}

type Service = platformbootstrap.RuntimeService[*gorm.DB, Dependencies, Runtime]
