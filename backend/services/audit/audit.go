package audit

import (
	"context"
	"errors"
	"log/slog"

	"gorm.io/gorm"
	"grovehub/libs/schoolgateway"
	"grovehub/services/audit/internal/app/bootstrap"
	bootstrapimpl "grovehub/services/audit/internal/app/bootstrap/impl"
	"grovehub/services/audit/internal/resources/migrations"
	"grovehub/services/infrastructure/client/messaging"
)

type Module = bootstrap.Runtime

func New(logger *slog.Logger, db *gorm.DB, events messaging.MessagingClient, access schoolgateway.ProtectedRouteAuthorizer) (Module, error) {
	return bootstrapimpl.New(logger, db, events, access).Build()
}

func Migrate(ctx context.Context, db *gorm.DB) error {
	if db == nil { return errors.New("audit migration requires a database") }
	sqlDB, err := db.DB()
	if err != nil { return err }
	return migrations.Migrate(ctx, sqlDB)
}
