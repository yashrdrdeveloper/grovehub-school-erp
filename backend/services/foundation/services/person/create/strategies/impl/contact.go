package impl

import (
    createhelpers "grovehub/services/foundation/services/person/create/helpers"
    "context"

    coreerrors "grovehub/libs/core/errors"
    servicepipeline "grovehub/libs/core/pipeline/service"
    personerror "grovehub/services/foundation/constants/error"
    persondto "grovehub/services/foundation/dto/person"
    "grovehub/services/foundation/entity"
    personflow "grovehub/services/foundation/services/person/create/flows"
)

func (strategy *Creator) Contacts(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person] {
    return servicepipeline.Decide(ctx, state,
        servicepipeline.When(
            func(state *personflow.Create) bool { return len(state.Contacts) == 0 },
            servicepipeline.Continue[personflow.Create, persondto.Person],
        ),
        servicepipeline.Otherwise(strategy.readContacts),
    )
}

func (strategy *Creator) readContacts(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person] {
    gids, err := strategy.contacts.FindPersonGIDsByContacts(ctx, state.Contacts)
    return servicepipeline.Result(gids, err,
        func(gids []string) servicepipeline.Outcome[persondto.Person] {
            return servicepipeline.Decide(ctx, state,
                servicepipeline.When(
                    func(*personflow.Create) bool { return len(gids) == 0 },
                    servicepipeline.Continue[personflow.Create, persondto.Person],
                ),
                servicepipeline.When(
                    func(*personflow.Create) bool { return len(gids) > 1 },
                    func(context.Context, *personflow.Create) servicepipeline.Outcome[persondto.Person] {
                        return servicepipeline.Fail[persondto.Person](coreerrors.NewError(personerror.PersonContactConflict))
                    },
                ),
                servicepipeline.Otherwise(func(ctx context.Context, _ *personflow.Create) servicepipeline.Outcome[persondto.Person] {
                    return strategy.returnPersonByGID(ctx, gids[0])
                }),
            )
        },
        createhelpers.FailRead,
    )
}

func (strategy *Creator) returnPersonByGID(ctx context.Context, gid string) servicepipeline.Outcome[persondto.Person] {
    person, err := strategy.persons.FindByGID(ctx, gid)
    return servicepipeline.Result(person, err,
        func(person entity.Person) servicepipeline.Outcome[persondto.Person] {
            return strategy.returnExisting(ctx, person)
        },
        createhelpers.FailRead,
    )
}
