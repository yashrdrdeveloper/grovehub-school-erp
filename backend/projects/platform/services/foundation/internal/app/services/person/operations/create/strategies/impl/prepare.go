package impl

import (
    createhelpers "grovehub/services/foundation/internal/app/services/person/operations/create/helpers"
    "context"

    coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
    coreidempotency "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency"
    headercontext "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext"
    corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    personflow "grovehub/services/foundation/internal/app/services/person/operations/create/flows"
    servicehelpers "grovehub/services/foundation/internal/app/services/person/helpers"
)

// Prepare builds domain contact values from an already validated request.
func (strategy *Creator) Prepare(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person] {
    state.Contacts = servicehelpers.BuildContacts(state.Request.RequestBody.Contacts, nil)
    return corepipeline.Decide(ctx, state,
        corepipeline.When(
            func(state *personflow.Create) bool { return len(state.Contacts) == 0 && headercontext.IdempotencyKey(ctx) == "" },
            func(context.Context, *personflow.Create) corepipeline.Outcome[persondto.Person] {
                return corepipeline.Fail[persondto.Person](coreerrors.NewError(personerror.InvalidPersonIdempotencyKey))
            },
        ),
        corepipeline.When(
            func(*personflow.Create) bool { return headercontext.IdempotencyKey(ctx) == "" },
            corepipeline.Continue[personflow.Create, persondto.Person],
        ),
        corepipeline.Otherwise(strategy.prepareIdempotency),
    )
}

func (strategy *Creator) prepareIdempotency(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person] {
    key, err := coreidempotency.ValidateKey(headercontext.IdempotencyKey(ctx))
    return corepipeline.ResolveResult(key, err,
        func(key string) corepipeline.Outcome[persondto.Person] {
            state.IdempotencyKey = key
            return strategy.fingerprint(state)
        },
        func(error) corepipeline.Outcome[persondto.Person] {
            return corepipeline.Fail[persondto.Person](coreerrors.NewError(personerror.InvalidPersonIdempotencyKey))
        },
    )
}

func (strategy *Creator) fingerprint(state *personflow.Create) corepipeline.Outcome[persondto.Person] {
    hash, err := coreidempotency.Fingerprint(state.Request.RequestBody)
    return corepipeline.ResolveResult(hash, err,
        func(hash string) corepipeline.Outcome[persondto.Person] {
            state.RequestHash = hash
            return corepipeline.Next[persondto.Person]()
        },
        createhelpers.FailCreate,
    )
}
