// Phase 1 annotation inventory. Each domain owns only its own declarations.
// The library remains the single metadata package; runtime capabilities stay
// in their existing owning core/service packages.
package annotation

func buildInventory() []UseCase {
    groups := [][]UseCase{
        entityUseCases,
        mappingUseCases,
        validationUseCases,
        persistenceUseCases,
        compositionUseCases,
        httpUseCases,
        middlewareUseCases,
        securityUseCases,
        externalUseCases,
        eventsUseCases,
        messagingUseCases,
        configurationUseCases,
        observabilityUseCases,
        serverUseCases,
        testingUseCases,
        utilitiesUseCases,
    }
    var result []UseCase
    for _, group := range groups { result = append(result, group...) }
    return result
}

// inventory contains every documented use case; name collisions intentionally
// represent one directive applied in different architectural contexts.
var inventory = buildInventory()
