// Package-level inventory composition. Each annotation category owns its
// definitions under its own directory. Keep public All/Lookup/Uses unchanged.
package annotation

import (
    "grovehub/libs/core/annotation/entity"
    "grovehub/libs/core/annotation/dto"
    "grovehub/libs/core/annotation/mapping"
    "grovehub/libs/core/annotation/validation"
    "grovehub/libs/core/annotation/persistence"
    "grovehub/libs/core/annotation/composition"
    "grovehub/libs/core/annotation/chain"
    "grovehub/libs/core/annotation/strategy"
    "grovehub/libs/core/annotation/http"
    "grovehub/libs/core/annotation/middleware"
    "grovehub/libs/core/annotation/security"
    "grovehub/libs/core/annotation/external"
    "grovehub/libs/core/annotation/event"
    "grovehub/libs/core/annotation/messaging"
    "grovehub/libs/core/annotation/configuration"
    "grovehub/libs/core/annotation/observability"
    "grovehub/libs/core/annotation/server"
    "grovehub/libs/core/annotation/testing"
    "grovehub/libs/core/annotation/utilities"
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
