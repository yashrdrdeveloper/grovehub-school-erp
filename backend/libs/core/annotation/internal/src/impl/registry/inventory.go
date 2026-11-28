// Package-level inventory composition. Each annotation category owns its
// definitions under its own directory. Keep public All/Lookup/Uses unchanged.
package registry

import (
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/entity"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/dto"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/mapping"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/validation"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/persistence"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/composition"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/chain"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/strategy"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/http"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/middleware"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/security"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/external"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/event"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/messaging"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/configuration"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/observability"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/server"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/testing"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/utilities"
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
