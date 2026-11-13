// Package create implements Person creation as independently testable strategies.
// Validation belongs to the controller, not this execution chain.
package create

import (
    "context"

    coreerrors "grovehub/libs/core/errors"
    coretransaction "grovehub/libs/core/database/transaction"
    servicepipeline "grovehub/libs/core/pipeline/service"
    personerror "grovehub/services/foundation/constants/error"
    persondto "grovehub/services/foundation/dto/person"
    "grovehub/services/foundation/entity"
    foundationmappers "grovehub/services/foundation/mappers/person"
    personrepository "grovehub/services/foundation/repository/person"
    contactrepository "grovehub/services/foundation/repository/contact"
    personflow "grovehub/services/foundation/services/person/flows"
)

// @Strategy
type Creator struct {
    persons personrepository.PersonRepository
    contacts contactrepository.ContactRepository
    unitOfWork coretransaction.UnitOfWork
    mapper foundationmappers.PersonMapper
}

func New(persons personrepository.PersonRepository, contacts contactrepository.ContactRepository, unit coretransaction.UnitOfWork, mapper foundationmappers.PersonMapper) *Creator {
    return &Creator{persons: persons, contacts: contacts, unitOfWork: unit, mapper: mapper}
}

func (strategy *Creator) returnExisting(ctx context.Context, person entity.Person) servicepipeline.Outcome[persondto.Person] {
    contacts, err := strategy.contacts.FindByPersonGID(ctx, person.GID)
    return servicepipeline.Result(contacts, err,
        func(stored []entity.Contact) servicepipeline.Outcome[persondto.Person] {
            return servicepipeline.Return(strategy.mapper.ToPersonResponse(person, stored), false)
        },
        func(err error) servicepipeline.Outcome[persondto.Person] {
            return servicepipeline.Fail[persondto.Person](coreerrors.Wrap(personerror.PersonReadFailed, err))
        },
    )
}

func failRead(err error) servicepipeline.Outcome[persondto.Person] {
    return servicepipeline.Fail[persondto.Person](coreerrors.Wrap(personerror.PersonReadFailed, err))
}

func failCreate(err error) servicepipeline.Outcome[persondto.Person] {
    return servicepipeline.Fail[persondto.Person](coreerrors.Wrap(personerror.PersonCreateFailed, err))
}

func (strategy *Creator) failCreateWithGID(_ context.Context, state *personflow.Create, err error) servicepipeline.Outcome[persondto.Person] {
    return servicepipeline.Fail[persondto.Person](
        coreerrors.WrapWithMetadata(personerror.PersonCreateFailed, err, strategy.mapper.ToPersonErrorMetadata(state.Person.GID)),
    )
}
