// Package foundation exposes only the module bootstrap entry point.
package foundation

import (
    "log/slog"

    "gorm.io/gorm"

    foundationbootstrap "grovehub/services/foundation/internal/app/bootstrap"
    bootstrapimpl "grovehub/services/foundation/internal/app/bootstrap/impl"
    infrastructureclient "grovehub/services/infrastructure/client"
)

// Module provides Foundation's shared lifecycle and typed Person client.
type Module = foundationbootstrap.Runtime

func New(logger *slog.Logger, db *gorm.DB, messaging infrastructureclient.MessagingClient) (Module, error) {
    return bootstrapimpl.New(logger, db, messaging).Build()
}
