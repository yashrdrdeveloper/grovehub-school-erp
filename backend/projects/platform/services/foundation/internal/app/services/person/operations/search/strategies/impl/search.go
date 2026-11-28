package impl

import (
    "context"
    "errors"

    coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
    coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
    servicepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/service"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    foundationmappers "grovehub/services/foundation/internal/app/mappers/person"
    personrepository "grovehub/services/foundation/internal/app/repository/person"
    contactrepository "grovehub/services/foundation/internal/app/repository/contact"
    searchflow "grovehub/services/foundation/internal/app/services/person/operations/search/flows"
    searchhelpers "grovehub/services/foundation/internal/app/services/person/operations/search/helpers"
    searchstrategies "grovehub/services/foundation/internal/app/services/person/operations/search/strategies"
    validatorhelpers "grovehub/services/foundation/internal/app/validators/person/helpers"
)

// @Strategy
type Searcher struct {
    persons personrepository.PersonRepository
    contacts contactrepository.ContactRepository
    mapper foundationmappers.PersonMapper
}

var _ searchstrategies.SearchStrategy = (*Searcher)(nil)

func New(persons personrepository.PersonRepository, contacts contactrepository.ContactRepository, mapper foundationmappers.PersonMapper) searchstrategies.SearchStrategy {
    return &Searcher{persons: persons, contacts: contacts, mapper: mapper}
}

func (strategy *Searcher) Prepare(ctx context.Context, state *searchflow.Search) servicepipeline.Outcome[coresearch.Result[persondto.Person]] {
    query, cursor, err := validatorhelpers.NormalizePersonSearchRequest(state.Request.RequestBody)
    if err != nil { return servicepipeline.Fail[coresearch.Result[persondto.Person]](err) }
    state.Query, state.Cursor = query, cursor
    return servicepipeline.Next[coresearch.Result[persondto.Person]]()
}

func (strategy *Searcher) Query(ctx context.Context, state *searchflow.Search) servicepipeline.Outcome[coresearch.Result[persondto.Person]] {
    result, err := strategy.persons.Search(ctx, state.Query)
    if errors.Is(err, coresearch.ErrInvalidQuery) {
        return servicepipeline.Fail[coresearch.Result[persondto.Person]](coreerrors.NewError(personerror.InvalidPersonSearch))
    }
    if err != nil { return servicepipeline.Fail[coresearch.Result[persondto.Person]](coreerrors.Wrap(personerror.PersonSearchFailed, err)) }
    state.Result = result
    return servicepipeline.Next[coresearch.Result[persondto.Person]]()
}

func (strategy *Searcher) Hydrate(ctx context.Context, state *searchflow.Search) servicepipeline.Outcome[coresearch.Result[persondto.Person]] {
    gids := make([]string, 0, len(state.Result.Items))
    for _, person := range state.Result.Items { gids = append(gids, person.GID) }
    contactsByGID, err := strategy.contacts.FindByPersonGIDs(ctx, gids)
    if err != nil { return servicepipeline.Fail[coresearch.Result[persondto.Person]](coreerrors.Wrap(personerror.PersonSearchFailed, err)) }
    items := make([]persondto.Person, 0, len(state.Result.Items))
    for _, person := range state.Result.Items {
        items = append(items, strategy.mapper.ToPersonResponse(person, contactsByGID[person.GID]))
    }
    nextCursor, err := searchhelpers.BuildNextCursor(state.Result, state.Query, state.Cursor)
    if err != nil { return servicepipeline.Fail[coresearch.Result[persondto.Person]](coreerrors.Wrap(personerror.PersonSearchFailed, err)) }
    pageInfo := coresearch.NewPageResponse(state.Query, state.Result)
    pageInfo.Rows = state.Query.Limit
    pageInfo.NextCursor = nextCursor
    return servicepipeline.Return(coresearch.Result[persondto.Person]{
        Items: items,
        PageInfo: pageInfo,
    }, false)
}
