// Package impl composes the Foundation Person operations into one service.
package impl

import (
    coretransaction "grovehub/libs/core/database/transaction"
    personrepository "grovehub/services/foundation/internal/app/repository/person"
    contactrepository "grovehub/services/foundation/internal/app/repository/contact"
    foundationmappers "grovehub/services/foundation/internal/app/mappers/person"
    foundationservices "grovehub/services/foundation/internal/app/services/person"
    foundationvalidators "grovehub/services/foundation/internal/app/validators/person"

    createoperation "grovehub/services/foundation/internal/app/services/person/operations/create"
    createimpl "grovehub/services/foundation/internal/app/services/person/operations/create/strategies/impl"
    getoperation "grovehub/services/foundation/internal/app/services/person/operations/getbygid"
    getimpl "grovehub/services/foundation/internal/app/services/person/operations/getbygid/strategies/impl"
    updateoperation "grovehub/services/foundation/internal/app/services/person/operations/update"
    updateimpl "grovehub/services/foundation/internal/app/services/person/operations/update/strategies/impl"
    searchoperation "grovehub/services/foundation/internal/app/services/person/operations/search"
    searchimpl "grovehub/services/foundation/internal/app/services/person/operations/search/strategies/impl"
    deactivateoperation "grovehub/services/foundation/internal/app/services/person/operations/deactivate"
    deactivateimpl "grovehub/services/foundation/internal/app/services/person/operations/deactivate/strategies/impl"
    eventproducers "grovehub/services/foundation/internal/app/events/producers/person"
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

func NewPersonService(personRepository personrepository.PersonRepository, contactRepository contactrepository.ContactRepository, unitOfWork coretransaction.UnitOfWork, mapper foundationmappers.PersonMapper, validator foundationvalidators.PersonValidator, publisher eventproducers.ChangedProducer) foundationservices.PersonService {
    if publisher == nil { panic("foundation PersonService requires an event producer") }
    return &PersonService{
        CreateOperation: createoperation.New(createimpl.NewCreateStrategy(personRepository, contactRepository, unitOfWork, mapper, publisher)),
        GetByGIDOperation: getoperation.New(getimpl.New(personRepository, contactRepository, mapper)),
        UpdateOperation: updateoperation.New(updateimpl.New(personRepository, contactRepository, unitOfWork, mapper, validator, publisher)),
        SearchOperation: searchoperation.New(searchimpl.New(personRepository, contactRepository, mapper)),
        DeactivateOperation: deactivateoperation.New(deactivateimpl.New(personRepository, contactRepository, unitOfWork, mapper, publisher)),
    }
}
