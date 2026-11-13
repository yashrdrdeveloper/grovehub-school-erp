package create

import (
    "context"
    "time"

    coreerrors "grovehub/libs/core/errors"
    coreid "grovehub/libs/core/id"
    servicepipeline "grovehub/libs/core/pipeline/service"
    personerror "grovehub/services/foundation/constants/error"
    persondto "grovehub/services/foundation/dto/person"
    repositoryerrors "grovehub/services/foundation/repository/constants/errors"
    personflow "grovehub/services/foundation/services/person/flows"
    servicehelpers "grovehub/services/foundation/services/person/helpers"
)

func (strategy *Creator) Persist(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person] {
    gid, err := coreid.NewUUIDv7()
    return servicepipeline.Result(gid, err,
        func(gid string) servicepipeline.Outcome[persondto.Person] {
            strategy.buildPerson(state, gid)
            return strategy.persistPrepared(ctx, state)
        },
        failCreate,
    )
}

func (strategy *Creator) buildPerson(state *personflow.Create, gid string) {
    person := servicehelpers.NewPerson(gid, state.Request.RequestBody, time.Now().UTC())
    servicepipeline.Decide(context.Background(), state,
        servicepipeline.When(
            func(state *personflow.Create) bool { return state.IdempotencyKey != "" },
            func(context.Context, *personflow.Create) servicepipeline.Outcome[persondto.Person] {
                person.IdempotencyKey = &state.IdempotencyKey
                person.RequestHash = &state.RequestHash
                return servicepipeline.Next[persondto.Person]()
            },
        ),
        servicepipeline.Otherwise(servicepipeline.Continue[personflow.Create, persondto.Person]),
    )
    state.Contacts = servicehelpers.AssociateContacts(person, state.Contacts)
    person.SearchDetails = servicehelpers.BuildPersonSearchDetails(person, state.Contacts)
    state.Person = person
}

func (strategy *Creator) persistPrepared(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person] {
    result := servicepipeline.Transactional(strategy.unitOfWork, strategy.persistWrites)(ctx, state)
    return servicepipeline.Recover(ctx, state, result,
        servicepipeline.OnError(repositoryerrors.ErrPersonIdempotencyKeyAlreadyExists, strategy.recoverIdempotency),
        servicepipeline.OnError(repositoryerrors.ErrContactAlreadyExists, strategy.recoverContacts),
        servicepipeline.OtherwiseError(strategy.failCreateWithGID),
    )
}

// @Transactional
func (strategy *Creator) persistWrites(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person] {
    return servicepipeline.Sequence(ctx, strategy.mapper.ToPersonResponse(state.Person, state.Contacts), true,
        func(ctx context.Context) error { return strategy.persons.Create(ctx, state.Person) },
        func(ctx context.Context) error { return strategy.contacts.ReplaceForPerson(ctx, state.Person.GID, state.Contacts) },
    )
}

func (strategy *Creator) recoverIdempotency(ctx context.Context, state *personflow.Create, _ error) servicepipeline.Outcome[persondto.Person] {
    return servicepipeline.RequireHandled(strategy.Idempotency(ctx, state), coreerrors.NewError(personerror.PersonReadFailed))
}

func (strategy *Creator) recoverContacts(ctx context.Context, state *personflow.Create, _ error) servicepipeline.Outcome[persondto.Person] {
    return servicepipeline.RequireHandled(strategy.Contacts(ctx, state), coreerrors.NewError(personerror.PersonCreateFailed))
}
