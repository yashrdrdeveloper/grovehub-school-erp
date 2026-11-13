package impl

import (
    "context"

    persondto "grovehub/services/foundation/dto/person"
    foundationmappers "grovehub/services/foundation/mappers/person"
    personrepository "grovehub/services/foundation/repository/person"
    contactrepository "grovehub/services/foundation/repository/contact"
    personhelpers "grovehub/services/foundation/services/person/helpers"
    readstrategies "grovehub/services/foundation/services/person/operations/getbygid/strategies"
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
