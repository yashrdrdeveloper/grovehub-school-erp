package impl

import (
    createhelpers "grovehub/services/foundation/internal/app/services/person/operations/create/helpers"
    "context"
    "time"

    coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
    coreid "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/id"
    corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    foundationevents "grovehub/services/foundation/internal/app/constants/events"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    repositoryerrors "grovehub/services/foundation/internal/app/repository/constants/errors"
    personflow "grovehub/services/foundation/internal/app/services/person/operations/create/flows"
    servicehelpers "grovehub/services/foundation/internal/app/services/person/helpers"
)

func (strategy *Creator) Persist(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person] {
    gid, err := coreid.NewUUIDv7()
    return corepipeline.ResolveResult(gid, err,
        func(gid string) corepipeline.Outcome[persondto.Person] {
            strategy.buildPerson(state, gid)
            return strategy.persistPrepared(ctx, state)
        },
        createhelpers.FailCreate,
    )
}

func (strategy *Creator) buildPerson(state *personflow.Create, gid string) {
    person := createhelpers.NewPerson(gid, state.Request.RequestBody, time.Now().UTC())
    corepipeline.Decide(context.Background(), state,
        corepipeline.When(
            func(state *personflow.Create) bool { return state.IdempotencyKey != "" },
            func(context.Context, *personflow.Create) corepipeline.Outcome[persondto.Person] {
                person.IdempotencyKey = &state.IdempotencyKey
                person.RequestHash = &state.RequestHash
                return corepipeline.Next[persondto.Person]()
            },
        ),
        corepipeline.Otherwise(corepipeline.Continue[personflow.Create, persondto.Person]),
    )
    state.Contacts = servicehelpers.AssociateContacts(person, state.Contacts)
    person.SearchDetails = servicehelpers.BuildPersonSearchDetails(person, state.Contacts)
    state.Person = person
}

func (strategy *Creator) persistPrepared(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person] {
    result := corepipeline.Transactional(strategy.unitOfWork, strategy.persistWrites)(ctx, state)
    return corepipeline.Recover(ctx, state, result,
        corepipeline.OnError(repositoryerrors.ErrPersonIdempotencyKeyAlreadyExists, strategy.recoverIdempotency),
        corepipeline.OnError(repositoryerrors.ErrContactAlreadyExists, strategy.recoverContacts),
        corepipeline.OtherwiseError(strategy.failCreateWithGID),
    )
}

// @Transactional
func (strategy *Creator) persistWrites(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person] {
    return corepipeline.Sequence(ctx, strategy.mapper.ToPersonResponse(state.Person, state.Contacts), true,
        func(ctx context.Context) error { return strategy.persons.Create(ctx, state.Person) },
        func(ctx context.Context) error { return strategy.contacts.ReplaceForPerson(ctx, state.Person.GID, state.Contacts) },
        func(ctx context.Context) error {
            next := strategy.mapper.ToPersonResponse(state.Person, state.Contacts)
            return strategy.events.Enqueue(ctx, foundationevents.Created, nil, &next)
        },
    )
}

func (strategy *Creator) recoverIdempotency(ctx context.Context, state *personflow.Create, _ error) corepipeline.Outcome[persondto.Person] {
    return corepipeline.RequireHandled(strategy.Idempotency(ctx, state), coreerrors.NewError(personerror.PersonReadFailed))
}

func (strategy *Creator) recoverContacts(ctx context.Context, state *personflow.Create, _ error) corepipeline.Outcome[persondto.Person] {
    return corepipeline.RequireHandled(strategy.Contacts(ctx, state), coreerrors.NewError(personerror.PersonCreateFailed))
}
