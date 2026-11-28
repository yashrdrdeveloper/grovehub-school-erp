package impl

import (
    "context"
    "errors"
    "time"

    corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
    coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
    foundationevents "grovehub/services/foundation/internal/app/constants/events"
    eventproducers "grovehub/services/foundation/internal/app/events/producers/person"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    "grovehub/services/foundation/internal/app/entity"
    foundationmappers "grovehub/services/foundation/internal/app/mappers/person"
    personrepository "grovehub/services/foundation/internal/app/repository/person"
    contactrepository "grovehub/services/foundation/internal/app/repository/contact"
    repositoryerrors "grovehub/services/foundation/internal/app/repository/constants/errors"
    personhelpers "grovehub/services/foundation/internal/app/services/person/helpers"
    deactivateflow "grovehub/services/foundation/internal/app/services/person/operations/deactivate/flows"
    deactivatestrategies "grovehub/services/foundation/internal/app/services/person/operations/deactivate/strategies"
)

// @Strategy
type Deactivator struct {
    persons personrepository.PersonRepository
    contacts contactrepository.ContactRepository
    unitOfWork coredatabase.UnitOfWork
    mapper foundationmappers.PersonMapper
    events eventproducers.ChangedProducer
}

var _ deactivatestrategies.DeactivateStrategy = (*Deactivator)(nil)

func New(persons personrepository.PersonRepository, contacts contactrepository.ContactRepository, unit coredatabase.UnitOfWork, mapper foundationmappers.PersonMapper, events eventproducers.ChangedProducer) deactivatestrategies.DeactivateStrategy {
    return &Deactivator{persons: persons, contacts: contacts, unitOfWork: unit, mapper: mapper, events: events}
}

func (strategy *Deactivator) Load(ctx context.Context, state *deactivateflow.Deactivate) corepipeline.Outcome[persondto.Person] {
    state.GID = strategy.mapper.ToPersonGID(state.Request.Params)
    person, contacts, err := personhelpers.ReadPerson(ctx, state.GID, strategy.persons, strategy.contacts, strategy.mapper)
    if err != nil { return corepipeline.Fail[persondto.Person](err) }
    state.Person, state.Contacts = person, contacts
    return corepipeline.Next[persondto.Person]()
}

func (strategy *Deactivator) Prepare(ctx context.Context, state *deactivateflow.Deactivate) corepipeline.Outcome[persondto.Person] {
    if state.Person.Status == entity.PersonStatusDeactivated {
        return corepipeline.Return(strategy.mapper.ToPersonResponse(state.Person, state.Contacts), false)
    }
    updated := state.Person
    updated.Status = entity.PersonStatusDeactivated
    updated.Version = state.Person.Version + 1
    updated.UpdatedAt = time.Now().UTC()
    state.UpdatedPerson = updated
    return corepipeline.Next[persondto.Person]()
}

// @Transactional
func (strategy *Deactivator) Persist(ctx context.Context, state *deactivateflow.Deactivate) corepipeline.Outcome[persondto.Person] {
    err := strategy.unitOfWork.Run(ctx, func(txCtx context.Context) error {
        if err := strategy.persons.Update(txCtx, state.UpdatedPerson); err != nil { return err }
        old := strategy.mapper.ToPersonResponse(state.Person, state.Contacts)
        next := strategy.mapper.ToPersonResponse(state.UpdatedPerson, state.Contacts)
        return strategy.events.Enqueue(txCtx, foundationevents.Deactivated, &old, &next)
    })
    if err != nil {
        if errors.Is(err, repositoryerrors.ErrPersonVersionConflict) {
            latest, contacts, readErr := personhelpers.ReadPerson(ctx, state.GID, strategy.persons, strategy.contacts, strategy.mapper)
            if readErr != nil { return corepipeline.Fail[persondto.Person](readErr) }
            if latest.Status == entity.PersonStatusDeactivated { return corepipeline.Return(strategy.mapper.ToPersonResponse(latest, contacts), false) }
        }
        return corepipeline.Fail[persondto.Person](personhelpers.MapUpdateError(state.GID, err, strategy.mapper))
    }
    return corepipeline.Return(strategy.mapper.ToPersonResponse(state.UpdatedPerson, state.Contacts), false)
}
