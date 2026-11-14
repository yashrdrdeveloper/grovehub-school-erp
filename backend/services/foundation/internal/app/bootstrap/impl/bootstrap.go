// Package impl composes Foundation's module dependencies behind Bootstrap.
package impl

import (
        "errors"
    "log/slog"
    "net/http"

    "gorm.io/gorm"

    coretelemetry "grovehub/libs/core/telemetry"
    foundationbootstrap "grovehub/services/foundation/internal/app/bootstrap"
    foundationclient "grovehub/services/foundation/client/person"
    protected "grovehub/services/foundation/internal/app/controller/person/protected"
    "grovehub/services/foundation/internal/app/routes"
    infrastructureclient "grovehub/services/infrastructure/client/messaging"
)

type builder struct {
    logger *slog.Logger
    db *gorm.DB
    events infrastructureclient.EventStore
}

var _ foundationbootstrap.Bootstrap = (*builder)(nil)

func New(logger *slog.Logger, db *gorm.DB, events infrastructureclient.EventStore) foundationbootstrap.Bootstrap {
    return &builder{logger: logger, db: db, events: events}
}

type runtime struct {
    logger *slog.Logger
    controller protected.PersonController
    client foundationclient.Client
}

var _ foundationbootstrap.Runtime = (*runtime)(nil)

func (bootstrap *builder) Build() (foundationbootstrap.Runtime, error) {
    if bootstrap.db == nil { return nil, errors.New("foundation bootstrap requires a database") }
    if bootstrap.events == nil { return nil, errors.New("foundation bootstrap requires an event store") }
    if bootstrap.logger == nil { bootstrap.logger = slog.Default() }

    producer := bootstrap.newEvents()
    controller := bootstrap.newRoutes(producer)
    client := bootstrap.newClient(controller)

    return &runtime{
        logger: bootstrap.logger,
        controller: controller,
        client: client,
    }, nil
}

func (service *runtime) RegisterHTTP(mux *http.ServeMux) {
    routes.RegisterHTTP(mux, service.controller)
    service.logger.Info(
        "foundation service registered",
        "service", coretelemetry.ServiceFoundation.Name,
        "serviceKey", coretelemetry.ServiceFoundation.Key,
    )
}

func (service *runtime) Client() foundationclient.Client { return service.client }

// Close remains part of the module contract; Infrastructure owns delivery.
func (service *runtime) Close() {}
