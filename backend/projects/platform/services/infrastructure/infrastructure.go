package infrastructure

import (
	platformbootstrap "grovehub/libs/platform/bootstrap"
    "context"
    "errors"
    "fmt"
    "log/slog"

    "gorm.io/gorm"

    messagingclient "grovehub/services/infrastructure/client/messaging"
    bootstrap "grovehub/services/infrastructure/internal/app/bootstrap"
    bootstrapimpl "grovehub/services/infrastructure/internal/app/bootstrap/impl"
    migrations "grovehub/services/infrastructure/internal/resources/migrations"
    repositoryimpl "grovehub/services/infrastructure/internal/app/repository/messaging/impl"
    mapperimpl "grovehub/services/infrastructure/internal/app/mappers/messaging/impl"
    validatorimpl "grovehub/services/infrastructure/internal/app/validators/messaging/impl"
)

type Module = bootstrap.Runtime

func New(config messagingclient.Config, db *gorm.DB, logger *slog.Logger) (Module, error) {
    return bootstrapimpl.New(config, db, logger).Build()
}

func Migrate(ctx context.Context, db *gorm.DB) error {
    if db == nil { return errors.New("infrastructure migrations require a database") }
    sqlDB, err := db.DB()
    if err != nil { return fmt.Errorf("access infrastructure migration database: %w", err) }
    return migrations.Migrate(ctx, sqlDB)
}

// NewEventStore provides the durable store without starting a broker dispatcher.
func NewEventStore(db *gorm.DB) messagingclient.EventStore {
    return repositoryimpl.New(db, validatorimpl.New(), mapperimpl.New())
}

// NewReplaySource exposes retained published events without broker or mutation access.
func NewReplaySource(db *gorm.DB) messagingclient.ReplaySource {
	return repositoryimpl.NewReplaySource(db)
}

// Bootstrap exposes operator capabilities without starting the service runtime.
func Bootstrap() platformbootstrap.Service[*gorm.DB] {
	return bootstrapimpl.NewOperatorBootstrap(Migrate)
}
