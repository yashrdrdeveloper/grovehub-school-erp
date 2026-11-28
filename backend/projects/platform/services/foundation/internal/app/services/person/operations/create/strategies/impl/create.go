// Package impl implements Person creation strategies.
// Validation belongs to the controller, not this execution chain.
package impl

import (
    "context"

    coreerrors "grovehub/libs/core/errors"
    coredatabase "grovehub/libs/core/database"
    createstrategies "grovehub/services/foundation/internal/app/services/person/operations/create/strategies"
    corepipeline "grovehub/libs/core/pipeline"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    "grovehub/services/foundation/internal/app/entity"
    foundationmappers "grovehub/services/foundation/internal/app/mappers/person"
    personrepository "grovehub/services/foundation/internal/app/repository/person"
    contactrepository "grovehub/services/foundation/internal/app/repository/contact"
    personflow "grovehub/services/foundation/internal/app/services/person/operations/create/flows"
    eventproducers "grovehub/services/foundation/internal/app/events/producers/person"
)

// @Strategy
type Creator struct {
    persons personrepository.PersonRepository
    contacts contactrepository.ContactRepository
    unitOfWork coredatabase.UnitOfWork
    mapper foundationmappers.PersonMapper
    events eventproducers.ChangedProducer
}

var _ createstrategies.CreateStrategy = (*Creator)(nil)

func NewCreateStrategy(persons personrepository.PersonRepository, contacts contactrepository.ContactRepository, unit coredatabase.UnitOfWork, mapper foundationmappers.PersonMapper, events eventproducers.ChangedProducer) createstrategies.CreateStrategy {
    return &Creator{persons: persons, contacts: contacts, unitOfWork: unit, mapper: mapper, events: events}
}

func (strategy *Creator) returnExisting(ctx context.Context, person entity.Person) corepipeline.Outcome[persondto.Person] {
    contacts, err := strategy.contacts.FindByPersonGID(ctx, person.GID)
    return corepipeline.ResolveResult(contacts, err,
        func(stored []entity.Contact) corepipeline.Outcome[persondto.Person] {
            return corepipeline.Return(strategy.mapper.ToPersonResponse(person, stored), false)
        },
        func(err error) corepipeline.Outcome[persondto.Person] {
            return corepipeline.Fail[persondto.Person](coreerrors.Wrap(personerror.PersonReadFailed, err))
        },
    )
}

func (strategy *Creator) failCreateWithGID(_ context.Context, state *personflow.Create, err error) corepipeline.Outcome[persondto.Person] {
    return corepipeline.Fail[persondto.Person](
        coreerrors.WrapWithMetadata(personerror.PersonCreateFailed, err, strategy.mapper.ToPersonErrorMetadata(state.Person.GID)),
    )
}
