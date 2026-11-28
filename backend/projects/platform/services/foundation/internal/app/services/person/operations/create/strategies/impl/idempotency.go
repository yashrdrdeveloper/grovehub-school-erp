package impl

import (
    createhelpers "grovehub/services/foundation/internal/app/services/person/operations/create/helpers"
    "context"

    coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
    servicepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/service"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    "grovehub/services/foundation/internal/app/entity"
    repositoryerrors "grovehub/services/foundation/internal/app/repository/constants/errors"
    personflow "grovehub/services/foundation/internal/app/services/person/operations/create/flows"
    coreutils "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils"
)

func (strategy *Creator) Idempotency(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person] {
    return servicepipeline.Decide(ctx, state,
        servicepipeline.When(
            func(state *personflow.Create) bool { return state.IdempotencyKey == "" },
            servicepipeline.Continue[personflow.Create, persondto.Person],
        ),
        servicepipeline.Otherwise(strategy.readIdempotency),
    )
}

func (strategy *Creator) readIdempotency(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person] {
    person, err := strategy.persons.FindByIdempotencyKey(ctx, state.IdempotencyKey)
    return servicepipeline.Lookup(person, err, repositoryerrors.ErrPersonNotFound,
        func(person entity.Person) servicepipeline.Outcome[persondto.Person] {
            return servicepipeline.Decide(ctx, state,
                servicepipeline.When(
                    func(*personflow.Create) bool { return coreutils.NullableString(person.RequestHash) != state.RequestHash },
                    func(context.Context, *personflow.Create) servicepipeline.Outcome[persondto.Person] {
                        return servicepipeline.Fail[persondto.Person](coreerrors.NewError(personerror.PersonIdempotencyConflict))
                    },
                ),
                servicepipeline.Otherwise(func(ctx context.Context, _ *personflow.Create) servicepipeline.Outcome[persondto.Person] {
                    return strategy.returnExisting(ctx, person)
                }),
            )
        },
        func() servicepipeline.Outcome[persondto.Person] { return servicepipeline.Next[persondto.Person]() },
        createhelpers.FailRead,
    )
}
