// Package impl composes the Foundation Person operations into one service.
package impl

import (
    coretransaction "grovehub/libs/core/database/transaction"
    personrepository "grovehub/services/foundation/repository/person"
    contactrepository "grovehub/services/foundation/repository/contact"
    foundationmappers "grovehub/services/foundation/mappers/person"
    foundationservices "grovehub/services/foundation/services/person"
    foundationvalidators "grovehub/services/foundation/validators/person"

    createoperation "grovehub/services/foundation/services/person/operations/create"
    createimpl "grovehub/services/foundation/services/person/operations/create/strategies/impl"
    getoperation "grovehub/services/foundation/services/person/operations/getbygid"
    getimpl "grovehub/services/foundation/services/person/operations/getbygid/strategies/impl"
    updateoperation "grovehub/services/foundation/services/person/operations/update"
    updateimpl "grovehub/services/foundation/services/person/operations/update/strategies/impl"
    searchoperation "grovehub/services/foundation/services/person/operations/search"
    searchimpl "grovehub/services/foundation/services/person/operations/search/strategies/impl"
    deactivateoperation "grovehub/services/foundation/services/person/operations/deactivate"
    deactivateimpl "grovehub/services/foundation/services/person/operations/deactivate/strategies/impl"
    eventproducers "grovehub/services/foundation/events/producers/person"
)

// @Service
type PersonService struct {
    *createoperation.CreateOperation
    *getoperation.GetByGIDOperation
    *updateoperation.UpdateOperation
    *searchoperation.SearchOperation
    *deactivateoperation.DeactivateOperation
}

var _ foundationservices.PersonService = (*PersonService)(nil)

func NewPersonService(personRepository personrepository.PersonRepository, contactRepository contactrepository.ContactRepository, unitOfWork coretransaction.UnitOfWork, mapper foundationmappers.PersonMapper, validator foundationvalidators.PersonValidator, publishers ...eventproducers.ChangedProducer) foundationservices.PersonService {
    publisher := eventproducers.Noop()
    if len(publishers) > 0 && publishers[0] != nil { publisher = publishers[0] }
    return &PersonService{
        CreateOperation: createoperation.New(createimpl.NewCreateStrategy(personRepository, contactRepository, unitOfWork, mapper, publisher)),
        GetByGIDOperation: getoperation.New(getimpl.New(personRepository, contactRepository, mapper)),
        UpdateOperation: updateoperation.New(updateimpl.New(personRepository, contactRepository, unitOfWork, mapper, validator, publisher)),
        SearchOperation: searchoperation.New(searchimpl.New(personRepository, contactRepository, mapper)),
        DeactivateOperation: deactivateoperation.New(deactivateimpl.New(personRepository, contactRepository, unitOfWork, mapper, publisher)),
    }
}
