package impl

import (
    createhelpers "grovehub/services/foundation/internal/app/services/person/operations/create/helpers"
    "context"

    coreerrors "grovehub/libs/core/errors"
    coreidempotency "grovehub/libs/core/idempotency"
    headercontext "grovehub/libs/core/headercontext"
    servicepipeline "grovehub/libs/core/pipeline/service"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    personflow "grovehub/services/foundation/internal/app/services/person/operations/create/flows"
    servicehelpers "grovehub/services/foundation/internal/app/services/person/helpers"
)

// Prepare builds domain contact values from an already validated request.
func (strategy *Creator) Prepare(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person] {
    state.Contacts = servicehelpers.BuildContacts(state.Request.RequestBody.Contacts, nil)
    return servicepipeline.Decide(ctx, state,
        servicepipeline.When(
            func(state *personflow.Create) bool { return len(state.Contacts) == 0 && headercontext.IdempotencyKey(ctx) == "" },
            func(context.Context, *personflow.Create) servicepipeline.Outcome[persondto.Person] {
                return servicepipeline.Fail[persondto.Person](coreerrors.NewError(personerror.InvalidPersonIdempotencyKey))
            },
        ),
        servicepipeline.When(
            func(*personflow.Create) bool { return headercontext.IdempotencyKey(ctx) == "" },
            servicepipeline.Continue[personflow.Create, persondto.Person],
        ),
        servicepipeline.Otherwise(strategy.prepareIdempotency),
    )
}

func (strategy *Creator) prepareIdempotency(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person] {
    key, err := coreidempotency.ValidateKey(headercontext.IdempotencyKey(ctx))
    return servicepipeline.Result(key, err,
        func(key string) servicepipeline.Outcome[persondto.Person] {
            state.IdempotencyKey = key
            return strategy.fingerprint(state)
        },
        func(error) servicepipeline.Outcome[persondto.Person] {
            return servicepipeline.Fail[persondto.Person](coreerrors.NewError(personerror.InvalidPersonIdempotencyKey))
        },
    )
}

func (strategy *Creator) fingerprint(state *personflow.Create) servicepipeline.Outcome[persondto.Person] {
    hash, err := coreidempotency.Fingerprint(state.Request.RequestBody)
    return servicepipeline.Result(hash, err,
        func(hash string) servicepipeline.Outcome[persondto.Person] {
            state.RequestHash = hash
            return servicepipeline.Next[persondto.Person]()
        },
        createhelpers.FailCreate,
    )
}
