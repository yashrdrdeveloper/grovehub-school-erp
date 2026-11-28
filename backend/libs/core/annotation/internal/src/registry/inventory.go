// Package-level inventory composition. Each annotation category owns its
// definitions under its own directory. Keep public All/Lookup/Uses unchanged.
package registry

import (
    "grovehub/libs/core/annotation/internal/src/catalog/entity"
    "grovehub/libs/core/annotation/internal/src/catalog/dto"
    "grovehub/libs/core/annotation/internal/src/catalog/mapping"
    "grovehub/libs/core/annotation/internal/src/catalog/validation"
    "grovehub/libs/core/annotation/internal/src/catalog/persistence"
    "grovehub/libs/core/annotation/internal/src/catalog/composition"
    "grovehub/libs/core/annotation/internal/src/catalog/chain"
    "grovehub/libs/core/annotation/internal/src/catalog/strategy"
    "grovehub/libs/core/annotation/internal/src/catalog/http"
    "grovehub/libs/core/annotation/internal/src/catalog/middleware"
    "grovehub/libs/core/annotation/internal/src/catalog/security"
    "grovehub/libs/core/annotation/internal/src/catalog/external"
    "grovehub/libs/core/annotation/internal/src/catalog/event"
    "grovehub/libs/core/annotation/internal/src/catalog/messaging"
    "grovehub/libs/core/annotation/internal/src/catalog/configuration"
    "grovehub/libs/core/annotation/internal/src/catalog/observability"
    "grovehub/libs/core/annotation/internal/src/catalog/server"
    "grovehub/libs/core/annotation/internal/src/catalog/testing"
    "grovehub/libs/core/annotation/internal/src/catalog/utilities"
)

func buildInventory() []UseCase {
    categories := [][]UseCase{
        entity.UseCases,
        dto.UseCases,
        mapping.UseCases,
        validation.UseCases,
        persistence.UseCases,
        composition.UseCases,
        chain.UseCases,
        strategy.UseCases,
        http.UseCases,
        middleware.UseCases,
        security.UseCases,
        external.UseCases,
        event.UseCases,
        messaging.UseCases,
        configuration.UseCases,
        observability.UseCases,
        server.UseCases,
        testing.UseCases,
        utilities.UseCases,
    }
    var definitions []UseCase
    for _, entries := range categories { definitions = append(definitions, entries...) }
    return definitions
}

var inventory = buildInventory()
