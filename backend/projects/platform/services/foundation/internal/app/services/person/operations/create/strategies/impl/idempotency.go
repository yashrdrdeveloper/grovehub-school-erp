package impl

import (
    createhelpers "grovehub/services/foundation/internal/app/services/person/operations/create/helpers"
    "context"

    coreerrors "grovehub/libs/core/errors"
    corepipeline "grovehub/libs/core/pipeline"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    "grovehub/services/foundation/internal/app/entity"
    repositoryerrors "grovehub/services/foundation/internal/app/repository/constants/errors"
    personflow "grovehub/services/foundation/internal/app/services/person/operations/create/flows"
    coreutils "grovehub/libs/core/utils"
)

func (strategy *Creator) Idempotency(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person] {
    return corepipeline.Decide(ctx, state,
        corepipeline.When(
            func(state *personflow.Create) bool { return state.IdempotencyKey == "" },
            corepipeline.Continue[personflow.Create, persondto.Person],
        ),
        corepipeline.Otherwise(strategy.readIdempotency),
    )
}

func (strategy *Creator) readIdempotency(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person] {
    person, err := strategy.persons.FindByIdempotencyKey(ctx, state.IdempotencyKey)
    return corepipeline.Lookup(person, err, repositoryerrors.ErrPersonNotFound,
        func(person entity.Person) corepipeline.Outcome[persondto.Person] {
            return corepipeline.Decide(ctx, state,
                corepipeline.When(
                    func(*personflow.Create) bool { return coreutils.NullableString(person.RequestHash) != state.RequestHash },
                    func(context.Context, *personflow.Create) corepipeline.Outcome[persondto.Person] {
                        return corepipeline.Fail[persondto.Person](coreerrors.NewError(personerror.PersonIdempotencyConflict))
                    },
                ),
                corepipeline.Otherwise(func(ctx context.Context, _ *personflow.Create) corepipeline.Outcome[persondto.Person] {
                    return strategy.returnExisting(ctx, person)
                }),
            )
        },
        func() corepipeline.Outcome[persondto.Person] { return corepipeline.Next[persondto.Person]() },
        createhelpers.FailRead,
    )
}
