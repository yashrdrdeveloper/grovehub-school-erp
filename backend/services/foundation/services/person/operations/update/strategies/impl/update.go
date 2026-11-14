package impl

import (
    "context"
    "errors"
    "time"

    annotationvalidation "grovehub/libs/core/annotation/validation"
    servicepipeline "grovehub/libs/core/pipeline/service"
    coretransaction "grovehub/libs/core/database/transaction"
    persondto "grovehub/services/foundation/dto/person"
    foundationevents "grovehub/services/foundation/constants/events"
    eventproducers "grovehub/services/foundation/events/person/producers"
    repositoryerrors "grovehub/services/foundation/repository/constants/errors"
    foundationmappers "grovehub/services/foundation/mappers/person"
    personrepository "grovehub/services/foundation/repository/person"
    contactrepository "grovehub/services/foundation/repository/contact"
    personhelpers "grovehub/services/foundation/services/person/helpers"
    updateflow "grovehub/services/foundation/services/person/operations/update/flows"
    updatehelpers "grovehub/services/foundation/services/person/operations/update/helpers"
    updatestrategies "grovehub/services/foundation/services/person/operations/update/strategies"
    foundationvalidators "grovehub/services/foundation/validators/person"
)

// @Strategy
type Updater struct {
    persons personrepository.PersonRepository
    contacts contactrepository.ContactRepository
    unitOfWork coretransaction.UnitOfWork
    mapper foundationmappers.PersonMapper
    validator foundationvalidators.PersonValidator
    events eventproducers.ChangedProducer
}

var _ updatestrategies.UpdateStrategy = (*Updater)(nil)

func New(persons personrepository.PersonRepository, contacts contactrepository.ContactRepository, unit coretransaction.UnitOfWork, mapper foundationmappers.PersonMapper, validator foundationvalidators.PersonValidator, events eventproducers.ChangedProducer) updatestrategies.UpdateStrategy {
    return &Updater{persons: persons, contacts: contacts, unitOfWork: unit, mapper: mapper, validator: validator, events: events}
}

func (strategy *Updater) Load(ctx context.Context, state *updateflow.Update) servicepipeline.Outcome[persondto.Person] {
    state.GID = strategy.mapper.ToPersonGID(state.Request.Params)
    person, contacts, err := personhelpers.ReadPerson(ctx, state.GID, strategy.persons, strategy.contacts, strategy.mapper)
    if err != nil { return servicepipeline.Fail[persondto.Person](err) }
    state.Person, state.Contacts = person, contacts
    return servicepipeline.Next[persondto.Person]()
}

func (strategy *Updater) Apply(ctx context.Context, state *updateflow.Update) servicepipeline.Outcome[persondto.Person] {
    validated, err := annotationvalidation.Process(state.Request, strategy.validator.UpdateRules())
    if err != nil { return servicepipeline.Fail[persondto.Person](err) }
    state.Validated = validated
    updated, contacts, changed, err := updatehelpers.ApplyPersonUpdate(ctx, state.Person, state.Contacts, validated.RequestBody, strategy.contacts, strategy.mapper)
    if err != nil { return servicepipeline.Fail[persondto.Person](err) }
    if !changed { return servicepipeline.Return(strategy.mapper.ToPersonResponse(state.Person, state.Contacts), false) }
    updated.Version = state.Person.Version + 1
    updated.UpdatedAt = time.Now().UTC()
    state.ContactsChanged = !updatehelpers.ContactsEqual(state.Contacts, contacts)
    if state.ContactsChanged { contacts = personhelpers.AssociateContacts(updated, contacts) }
    updated.SearchDetails = personhelpers.BuildPersonSearchDetails(updated, contacts)
    state.UpdatedPerson, state.UpdatedContacts = updated, contacts
    return servicepipeline.Next[persondto.Person]()
}

// @Transactional
func (strategy *Updater) Persist(ctx context.Context, state *updateflow.Update) servicepipeline.Outcome[persondto.Person] {
    err := strategy.unitOfWork.Run(ctx, func(txCtx context.Context) error {
        if err := strategy.persons.Update(txCtx, state.UpdatedPerson); err != nil { return err }
        if state.ContactsChanged {
            if err := strategy.contacts.ReplaceForPerson(txCtx, state.UpdatedPerson.GID, state.UpdatedContacts); err != nil { return err }
        }
        old := strategy.mapper.ToPersonResponse(state.Person, state.Contacts)
        next := strategy.mapper.ToPersonResponse(state.UpdatedPerson, state.UpdatedContacts)
        return strategy.events.Enqueue(txCtx, foundationevents.Updated, &old, &next)
    })
    if err != nil {
        if errors.Is(err, repositoryerrors.ErrPersonVersionConflict) {
            latest, latestContacts, readErr := personhelpers.ReadPerson(ctx, state.GID, strategy.persons, strategy.contacts, strategy.mapper)
            if readErr != nil { return servicepipeline.Fail[persondto.Person](readErr) }
            _, _, stillChanged, reconcileErr := updatehelpers.ApplyPersonUpdate(ctx, latest, latestContacts, state.Validated.RequestBody, strategy.contacts, strategy.mapper)
            if reconcileErr != nil { return servicepipeline.Fail[persondto.Person](reconcileErr) }
            if !stillChanged { return servicepipeline.Return(strategy.mapper.ToPersonResponse(latest, latestContacts), false) }
        }
        return servicepipeline.Fail[persondto.Person](personhelpers.MapUpdateError(state.GID, err, strategy.mapper))
    }
    return servicepipeline.Return(strategy.mapper.ToPersonResponse(state.UpdatedPerson, state.UpdatedContacts), false)
}
