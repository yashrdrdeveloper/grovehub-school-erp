package impl

import (
    "context"

    persondto "grovehub/services/foundation/client/person/dto"
    foundationmappers "grovehub/services/foundation/internal/app/mappers/person"
    personrepository "grovehub/services/foundation/internal/app/repository/person"
    contactrepository "grovehub/services/foundation/internal/app/repository/contact"
    personhelpers "grovehub/services/foundation/internal/app/services/person/helpers"
    readstrategies "grovehub/services/foundation/internal/app/services/person/operations/getbygid/strategies"
)

// @Strategy
type Getter struct {
    persons personrepository.PersonRepository
    contacts contactrepository.ContactRepository
    mapper foundationmappers.PersonMapper
}

var _ readstrategies.GetByGIDStrategy = (*Getter)(nil)

func New(persons personrepository.PersonRepository, contacts contactrepository.ContactRepository, mapper foundationmappers.PersonMapper) readstrategies.GetByGIDStrategy {
    return &Getter{persons: persons, contacts: contacts, mapper: mapper}
}

func (strategy *Getter) Execute(ctx context.Context, request persondto.GetPersonByGIDRequest) (persondto.Person, error) {
    person, contacts, err := personhelpers.ReadPerson(ctx, strategy.mapper.ToPersonGID(request.Params), strategy.persons, strategy.contacts, strategy.mapper)
    if err != nil { return persondto.Person{}, err }
    return strategy.mapper.ToPersonResponse(person, contacts), nil
}
