package impl

import (
    "context"
    "errors"
    "time"

    servicepipeline "grovehub/libs/core/pipeline/service"
    persondto "grovehub/services/foundation/dto/person"
    "grovehub/services/foundation/entity"
    foundationmappers "grovehub/services/foundation/mappers/person"
    personrepository "grovehub/services/foundation/repository/person"
    contactrepository "grovehub/services/foundation/repository/contact"
    repositoryerrors "grovehub/services/foundation/repository/constants/errors"
    personhelpers "grovehub/services/foundation/services/person/helpers"
    deactivateflow "grovehub/services/foundation/services/person/operations/deactivate/flows"
    deactivatestrategies "grovehub/services/foundation/services/person/operations/deactivate/strategies"
)

// @Strategy
type Deactivator struct {
    persons personrepository.PersonRepository
    contacts contactrepository.ContactRepository
    mapper foundationmappers.PersonMapper
}

var _ deactivatestrategies.DeactivateStrategy = (*Deactivator)(nil)

func New(persons personrepository.PersonRepository, contacts contactrepository.ContactRepository, mapper foundationmappers.PersonMapper) deactivatestrategies.DeactivateStrategy {
    return &Deactivator{persons: persons, contacts: contacts, mapper: mapper}
}

func (strategy *Deactivator) Load(ctx context.Context, state *deactivateflow.Deactivate) servicepipeline.Outcome[persondto.Person] {
    state.GID = strategy.mapper.ToPersonGID(state.Request.Params)
    person, contacts, err := personhelpers.ReadPerson(ctx, state.GID, strategy.persons, strategy.contacts, strategy.mapper)
    if err != nil { return servicepipeline.Fail[persondto.Person](err) }
    state.Person, state.Contacts = person, contacts
    return servicepipeline.Next[persondto.Person]()
}

func (strategy *Deactivator) Prepare(ctx context.Context, state *deactivateflow.Deactivate) servicepipeline.Outcome[persondto.Person] {
    if state.Person.Status == entity.PersonStatusDeactivated {
        return servicepipeline.Return(strategy.mapper.ToPersonResponse(state.Person, state.Contacts), false)
    }
    updated := state.Person
    updated.Status = entity.PersonStatusDeactivated
    updated.Version = state.Person.Version + 1
    updated.UpdatedAt = time.Now().UTC()
    state.UpdatedPerson = updated
    return servicepipeline.Next[persondto.Person]()
}

func (strategy *Deactivator) Persist(ctx context.Context, state *deactivateflow.Deactivate) servicepipeline.Outcome[persondto.Person] {
    if err := strategy.persons.Update(ctx, state.UpdatedPerson); err != nil {
        if errors.Is(err, repositoryerrors.ErrPersonVersionConflict) {
            latest, contacts, readErr := personhelpers.ReadPerson(ctx, state.GID, strategy.persons, strategy.contacts, strategy.mapper)
            if readErr != nil { return servicepipeline.Fail[persondto.Person](readErr) }
            if latest.Status == entity.PersonStatusDeactivated { return servicepipeline.Return(strategy.mapper.ToPersonResponse(latest, contacts), false) }
        }
        return servicepipeline.Fail[persondto.Person](personhelpers.MapUpdateError(state.GID, err, strategy.mapper))
    }
    return servicepipeline.Return(strategy.mapper.ToPersonResponse(state.UpdatedPerson, state.Contacts), false)
}
