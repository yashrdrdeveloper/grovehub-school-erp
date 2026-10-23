package impl

import (
	"context"

	searchcontract "grovehub/libs/core/contracts/search"
	foundationclient "grovehub/services/foundation/client"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
)

type personService interface {
	Create(
		ctx context.Context,
		request foundationrequest.CreatePerson,
	) (foundationresponse.Person, bool, error)
	GetByGID(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
	Update(
		ctx context.Context,
		personRequest foundationrequest.PersonByGID,
		request foundationrequest.UpdatePerson,
	) (foundationresponse.Person, error)
	Search(
		ctx context.Context,
		request searchcontract.Request,
	) (searchcontract.Result[foundationresponse.Person], error)
	Deactivate(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
}

type personClient struct {
	personService personService
}

func NewPersonClient(personService personService) foundationclient.Client {
	return &personClient{personService: personService}
}

func (client *personClient) CreatePerson(
	ctx context.Context,
	request foundationrequest.CreatePerson,
) (foundationresponse.Person, bool, error) {
	return client.personService.Create(ctx, request)
}

func (client *personClient) GetPersonByGID(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	return client.personService.GetByGID(ctx, request)
}

func (client *personClient) UpdatePerson(
	ctx context.Context,
	personRequest foundationrequest.PersonByGID,
	request foundationrequest.UpdatePerson,
) (foundationresponse.Person, error) {
	return client.personService.Update(ctx, personRequest, request)
}

func (client *personClient) SearchPersons(
	ctx context.Context,
	request searchcontract.Request,
) (searchcontract.Result[foundationresponse.Person], error) {
	return client.personService.Search(ctx, request)
}

func (client *personClient) DeactivatePerson(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	return client.personService.Deactivate(ctx, request)
}
