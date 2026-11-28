package impl

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	platformbootstrap "grovehub/libs/platform/bootstrap"
	"grovehub/services/integration/internal/app/bootstrap"
	"grovehub/services/integration/internal/resources/migrations"
	"net/http"
	"grovehub/services/identity/client/invitation"
	"grovehub/services/integration/client/mail"
	clientimpl "grovehub/services/integration/internal/app/client/mail/impl"
	controller "grovehub/services/integration/internal/app/controller/mail/restricted"
	service "grovehub/services/integration/internal/app/services/mail/impl"
)

// NewServiceBootstrap joins runtime, command and migration registration.
func NewServiceBootstrap() bootstrap.Service {
	return platformbootstrap.RuntimeRegistration[*gorm.DB, bootstrap.Dependencies, bootstrap.Runtime]{
		Service: NewOperatorBootstrap(runMigration),
		BuildRuntime: func(input bootstrap.Dependencies) (bootstrap.Runtime, error) {
			return New(input.Logger, input.DB, input.Events, input.DevAuth, input.TenantState...).Build()
		},
	}
}

func runMigration(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return errors.New("integration migrations require a database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("access integration migration database: %w", err)
	}
	return migrations.Migrate(ctx, sqlDB)
}

func NewMailReconciliationClient(db *gorm.DB, jobs invitation.DeliveryReader) mail.Client {
	return clientimpl.NewMailClient(service.New(db, jobs))
}

func RegisterMailReconciliationHTTP(mux *http.ServeMux, client mail.Client) {
	mux.HandleFunc("POST "+mail.ReconciliationPath, controller.SearchHandler(client))
}
