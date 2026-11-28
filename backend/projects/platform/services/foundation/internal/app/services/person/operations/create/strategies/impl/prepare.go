package impl

import (
	"context"
	createhelpers "grovehub/services/foundation/internal/app/services/person/operations/create/helpers"

	coreerrors "grovehub/libs/core/errors"
	headercontext "grovehub/libs/core/headercontext"
	coreidempotency "grovehub/libs/core/idempotency"
	corepipeline "grovehub/libs/core/pipeline"
	personerror "grovehub/services/foundation/internal/app/constants/error"
	persondto "grovehub/services/foundation/internal/app/dto/person"
	servicehelpers "grovehub/services/foundation/internal/app/services/person/helpers"
	personflow "grovehub/services/foundation/internal/app/services/person/operations/create/flows"
)

// Prepare builds domain contact values from an already validated request.
func (strategy *Creator) Prepare(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person] {
	state.Contacts = servicehelpers.BuildContacts(state.Request.RequestBody.Contacts, nil)
	return corepipeline.Decide(ctx, state,
		corepipeline.When(
			func(state *personflow.Create) bool {
				return len(state.Contacts) == 0 && headercontext.IdempotencyKey(ctx) == ""
			},
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
