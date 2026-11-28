package impl

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	platformbootstrap "grovehub/libs/platform/bootstrap"
	"grovehub/services/identity/internal/app/bootstrap"
	"grovehub/services/identity/internal/resources/migrations"
	"log/slog"
	"grovehub/libs/platform/session"
	"grovehub/services/identity/client/invitation"
	clientimpl "grovehub/services/identity/internal/app/client/invitation/impl"
	repository "grovehub/services/identity/internal/app/repository/invitation/impl"
	sessionrepo "grovehub/services/identity/internal/app/repository/session/impl"
	"grovehub/services/identity/internal/app/services/authentication/operations/reconciletenant"
)

// NewServiceBootstrap joins runtime, command and migration registration.
func NewServiceBootstrap() bootstrap.Service {
	return platformbootstrap.RuntimeRegistration[*gorm.DB, bootstrap.Dependencies, bootstrap.Runtime]{
		Service: NewOperatorBootstrap(runMigration),
		BuildRuntime: func(input bootstrap.Dependencies) (bootstrap.Runtime, error) {
			return New(input.Logger, input.DB, input.Foundation, input.Integration, input.Events, input.TenantState...).Build()
		},
	}
}

func runMigration(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("identity migration requires a database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("identity migration DB: %w", err)
	}
	return migrations.Migrate(ctx, sqlDB)
}

func NewInvitationDeliveryReader(db *gorm.DB) invitation.DeliveryReader {
	return clientimpl.NewInvitationClient(repository.NewDeliveryReader(db))
}

func StartTenantSessionReconciler(parent context.Context, logger *slog.Logger, db *gorm.DB, lookup session.TenantStateLookup) func() {
	if db == nil || lookup == nil {
		return func() {}
	}
	ctx, cancel := context.WithCancel(parent)
	done := make(chan struct{})
	worker := reconciletenant.New(sessionrepo.NewTenantSessions(db), lookup, logger)
	go func() {
		defer close(done)
		worker.Run(ctx)
	}()
	return func() {
		cancel()
		<-done
	}
}
