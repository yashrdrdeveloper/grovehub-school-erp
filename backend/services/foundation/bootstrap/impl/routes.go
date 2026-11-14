package impl

import (
    coretransaction "grovehub/libs/core/database/transaction"
    protected "grovehub/services/foundation/controller/person/protected"
    controllerimpl "grovehub/services/foundation/controller/person/protected/impl"
    personproducer "grovehub/services/foundation/events/producers/person"
    mapperimpl "grovehub/services/foundation/mappers/person/impl"
    contactrepositoryimpl "grovehub/services/foundation/repository/contact/impl"
    personrepositoryimpl "grovehub/services/foundation/repository/person/impl"
    serviceimpl "grovehub/services/foundation/services/person/impl"
    validatorimpl "grovehub/services/foundation/validators/person/impl"
)

// newRoutes composes Person dependencies for the Foundation HTTP routes.
// Actual URL registration remains in foundation/routes/http.go.
func (bootstrap *builder) newRoutes(producer personproducer.ChangedProducer) protected.PersonController {
    persons := personrepositoryimpl.NewPersonRepository(bootstrap.db)
    contacts := contactrepositoryimpl.NewContactRepository(bootstrap.db)
    unitOfWork := coretransaction.NewUnitOfWork(bootstrap.db)
    mapper := mapperimpl.NewPersonMapper()
    validator := validatorimpl.NewPersonValidator()
    service := serviceimpl.NewPersonService(persons, contacts, unitOfWork, mapper, validator, producer)
    return controllerimpl.NewPersonController(service, validator)
}
