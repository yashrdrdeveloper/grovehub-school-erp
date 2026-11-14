// Package impl composes Foundation's module dependencies behind Bootstrap.
package impl

import (
    "context"
    "errors"
    "log/slog"
    "net/http"

    "gorm.io/gorm"

    coretelemetry "grovehub/libs/core/telemetry"
    foundationbootstrap "grovehub/services/foundation/bootstrap"
    foundationclient "grovehub/services/foundation/client/person"
    protected "grovehub/services/foundation/controller/person/protected"
    "grovehub/services/foundation/routes"
    infrastructureclient "grovehub/services/infrastructure/client"
)

type builder struct {
    logger *slog.Logger
    db *gorm.DB
    messaging infrastructureclient.MessagingClient
}

var _ foundationbootstrap.Bootstrap = (*builder)(nil)

func New(logger *slog.Logger, db *gorm.DB, messaging infrastructureclient.MessagingClient) foundationbootstrap.Bootstrap {
    return &builder{logger: logger, db: db, messaging: messaging}
}

type runtime struct {
    logger *slog.Logger
    controller protected.PersonController
    client foundationclient.Client
    stopDispatcher context.CancelFunc
    dispatcherDone <-chan struct{}
}

var _ foundationbootstrap.Runtime = (*runtime)(nil)

func (bootstrap *builder) Build() (foundationbootstrap.Runtime, error) {
    if bootstrap.db == nil { return nil, errors.New("foundation bootstrap requires a database") }
    if bootstrap.messaging == nil { return nil, errors.New("foundation bootstrap requires a messaging client") }
    if bootstrap.logger == nil { bootstrap.logger = slog.Default() }

    events := bootstrap.newEvents()
    controller := bootstrap.newRoutes(events.producer)
    client := bootstrap.newClient(controller)
    cancel, done := events.start()

    return &runtime{
        logger: bootstrap.logger,
        controller: controller,
        client: client,
        stopDispatcher: cancel,
        dispatcherDone: done,
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

// Close stops Foundation-owned delivery before the database and broker close.
func (service *runtime) Close() {
    service.stopDispatcher()
    <-service.dispatcherDone
}
