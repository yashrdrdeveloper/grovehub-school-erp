package impl

import (
    createhelpers "grovehub/services/foundation/internal/app/services/person/operations/create/helpers"
    "context"

    coreerrors "grovehub/libs/core/errors"
    corepipeline "grovehub/libs/core/pipeline"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    "grovehub/services/foundation/internal/app/entity"
    personflow "grovehub/services/foundation/internal/app/services/person/operations/create/flows"
)

func (strategy *Creator) Contacts(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person] {
    return corepipeline.Decide(ctx, state,
        corepipeline.When(
            func(state *personflow.Create) bool { return len(state.Contacts) == 0 },
            corepipeline.Continue[personflow.Create, persondto.Person],
        ),
        corepipeline.Otherwise(strategy.readContacts),
    )
}

func (strategy *Creator) readContacts(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person] {
    gids, err := strategy.contacts.FindPersonGIDsByContacts(ctx, state.Contacts)
    return corepipeline.ResolveResult(gids, err,
        func(gids []string) corepipeline.Outcome[persondto.Person] {
            return corepipeline.Decide(ctx, state,
                corepipeline.When(
                    func(*personflow.Create) bool { return len(gids) == 0 },
                    corepipeline.Continue[personflow.Create, persondto.Person],
                ),
                corepipeline.When(
                    func(*personflow.Create) bool { return len(gids) > 1 },
                    func(context.Context, *personflow.Create) corepipeline.Outcome[persondto.Person] {
                        return corepipeline.Fail[persondto.Person](coreerrors.NewError(personerror.PersonContactConflict))
                    },
                ),
                corepipeline.Otherwise(func(ctx context.Context, _ *personflow.Create) corepipeline.Outcome[persondto.Person] {
                    return strategy.returnPersonByGID(ctx, gids[0])
                }),
            )
        },
        createhelpers.FailRead,
    )
}

func (strategy *Creator) returnPersonByGID(ctx context.Context, gid string) corepipeline.Outcome[persondto.Person] {
    person, err := strategy.persons.FindByGID(ctx, gid)
    return corepipeline.ResolveResult(person, err,
        func(person entity.Person) corepipeline.Outcome[persondto.Person] {
            return strategy.returnExisting(ctx, person)
        },
        createhelpers.FailRead,
    )
}
