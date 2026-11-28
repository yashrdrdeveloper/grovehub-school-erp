// Package integration exposes the service bootstrap and compatibility entry points.
package integration

import (
	"context"
	"log/slog"
	"gorm.io/gorm"
	"grovehub/libs/platform/session"
	"grovehub/services/infrastructure/client/messaging"
	"grovehub/services/integration/internal/app/bootstrap"
	bootstrapimpl "grovehub/services/integration/internal/app/bootstrap/impl"
	"net/http"
	"grovehub/services/identity/client/invitation"
	"grovehub/services/integration/client/mail"
)

type Module = bootstrap.Runtime
type Dependencies = bootstrap.Dependencies
type Service = bootstrap.Service

// Bootstrap is the common typed service entry point.
func Bootstrap() Service {
	return bootstrapimpl.NewServiceBootstrap()
}

// Migrate is a compatibility adapter for callers migrating to Bootstrap().Migration.
func Migrate(ctx context.Context, db *gorm.DB) error {
	return Bootstrap().Migration().Run(ctx, db)
}

type DevAuth = bootstrap.DevAuth

func New(logger *slog.Logger, db *gorm.DB, events messaging.MessagingClient, developmentAuth ...DevAuth) (Module, error) {
	settings := DevAuth{}
	if len(developmentAuth) > 0 {
		settings = developmentAuth[0]
	}
	return Bootstrap().Build(Dependencies{Logger: logger, DB: db, Events: events, DevAuth: settings})
}

func NewWithTenantState(logger *slog.Logger, db *gorm.DB, events messaging.MessagingClient, lookup session.TenantStateLookup) (Module, error) {
	return Bootstrap().Build(Dependencies{Logger: logger, DB: db, Events: events, TenantState: []session.TenantStateLookup{lookup}})
}

func NewMailReconciliationClient(db *gorm.DB, jobs invitation.DeliveryReader) mail.Client {
	return bootstrapimpl.NewMailReconciliationClient(db, jobs)
}

func RegisterMailReconciliationHTTP(mux *http.ServeMux, client mail.Client) {
	bootstrapimpl.RegisterMailReconciliationHTTP(mux, client)
}
