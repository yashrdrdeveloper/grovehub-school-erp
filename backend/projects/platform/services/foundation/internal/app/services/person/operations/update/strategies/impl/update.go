package impl

import (
	"context"
	"errors"
	"time"

	coreannotation "grovehub/libs/core/annotation"
	coredatabase "grovehub/libs/core/database"
	corepipeline "grovehub/libs/core/pipeline"
	foundationevents "grovehub/services/foundation/internal/app/constants/events"
	persondto "grovehub/services/foundation/internal/app/dto/person"
	eventproducers "grovehub/services/foundation/internal/app/events/producers/person"
	foundationmappers "grovehub/services/foundation/internal/app/mappers/person"
	repositoryerrors "grovehub/services/foundation/internal/app/repository/constants/errors"
	contactrepository "grovehub/services/foundation/internal/app/repository/contact"
	personrepository "grovehub/services/foundation/internal/app/repository/person"
	personhelpers "grovehub/services/foundation/internal/app/services/person/helpers"
	updateflow "grovehub/services/foundation/internal/app/services/person/operations/update/flows"
	updatehelpers "grovehub/services/foundation/internal/app/services/person/operations/update/helpers"
	updatestrategies "grovehub/services/foundation/internal/app/services/person/operations/update/strategies"
	foundationvalidators "grovehub/services/foundation/internal/app/validators/person"
)

// @Strategy
type Updater struct {
	persons    personrepository.PersonRepository
	contacts   contactrepository.ContactRepository
	unitOfWork coredatabase.UnitOfWork
	mapper     foundationmappers.PersonMapper
	validator  foundationvalidators.PersonValidator
	events     eventproducers.ChangedProducer
}

var _ updatestrategies.UpdateStrategy = (*Updater)(nil)

func New(persons personrepository.PersonRepository, contacts contactrepository.ContactRepository, unit coredatabase.UnitOfWork, mapper foundationmappers.PersonMapper, validator foundationvalidators.PersonValidator, events eventproducers.ChangedProducer) updatestrategies.UpdateStrategy {
	return &Updater{persons: persons, contacts: contacts, unitOfWork: unit, mapper: mapper, validator: validator, events: events}
}

func (strategy *Updater) Load(ctx context.Context, state *updateflow.Update) corepipeline.Outcome[persondto.Person] {
	state.GID = strategy.mapper.ToPersonGID(state.Request.Params)
	person, contacts, err := personhelpers.ReadPerson(ctx, state.GID, strategy.persons, strategy.contacts, strategy.mapper)
	if err != nil {
		return corepipeline.Fail[persondto.Person](err)
	}
	state.Person, state.Contacts = person, contacts
	return corepipeline.Next[persondto.Person]()
}

func (strategy *Updater) Apply(ctx context.Context, state *updateflow.Update) corepipeline.Outcome[persondto.Person] {
	validated, err := coreannotation.Process(state.Request, strategy.validator.UpdateRules())
	if err != nil {
		return corepipeline.Fail[persondto.Person](err)
	}
	state.Validated = validated
	updated, contacts, changed, err := updatehelpers.ApplyPersonUpdate(ctx, state.Person, state.Contacts, validated.RequestBody, strategy.contacts, strategy.mapper)
	if err != nil {
		return corepipeline.Fail[persondto.Person](err)
	}
	if !changed {
		return corepipeline.Return(strategy.mapper.ToPersonResponse(state.Person, state.Contacts), false)
	}
	updated.Version = state.Person.Version + 1
	updated.UpdatedAt = time.Now().UTC()
	state.ContactsChanged = !updatehelpers.ContactsEqual(state.Contacts, contacts)
	if state.ContactsChanged {
		contacts = personhelpers.AssociateContacts(updated, contacts)
	}
	updated.SearchDetails = personhelpers.BuildPersonSearchDetails(updated, contacts)
	state.UpdatedPerson, state.UpdatedContacts = updated, contacts
	return corepipeline.Next[persondto.Person]()
}

// @Transactional
func (strategy *Updater) Persist(ctx context.Context, state *updateflow.Update) corepipeline.Outcome[persondto.Person] {
	err := strategy.unitOfWork.Run(ctx, func(txCtx context.Context) error {
		if err := strategy.persons.Update(txCtx, state.UpdatedPerson); err != nil {
			return err
		}
		if state.ContactsChanged {
			if err := strategy.contacts.ReplaceForPerson(txCtx, state.UpdatedPerson.GID, state.UpdatedContacts); err != nil {
				return err
			}
		}
		old := strategy.mapper.ToPersonResponse(state.Person, state.Contacts)
		next := strategy.mapper.ToPersonResponse(state.UpdatedPerson, state.UpdatedContacts)
		return strategy.events.Enqueue(txCtx, foundationevents.Updated, &old, &next)
	})
	if err != nil {
		if errors.Is(err, repositoryerrors.ErrPersonVersionConflict) {
			latest, latestContacts, readErr := personhelpers.ReadPerson(ctx, state.GID, strategy.persons, strategy.contacts, strategy.mapper)
			if readErr != nil {
				return corepipeline.Fail[persondto.Person](readErr)
			}
			_, _, stillChanged, reconcileErr := updatehelpers.ApplyPersonUpdate(ctx, latest, latestContacts, state.Validated.RequestBody, strategy.contacts, strategy.mapper)
			if reconcileErr != nil {
				return corepipeline.Fail[persondto.Person](reconcileErr)
			}
			if !stillChanged {
				return corepipeline.Return(strategy.mapper.ToPersonResponse(latest, latestContacts), false)
			}
		}
		return corepipeline.Fail[persondto.Person](personhelpers.MapUpdateError(state.GID, err, strategy.mapper))
	}
	return corepipeline.Return(strategy.mapper.ToPersonResponse(state.UpdatedPerson, state.UpdatedContacts), false)
}
