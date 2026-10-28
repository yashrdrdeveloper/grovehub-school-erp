package impl

import (
	"context"

	coresearch "grovehub/libs/core/search"
	foundationclient "grovehub/services/foundation/client"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	foundationservices "grovehub/services/foundation/services"
)

type personClient struct {
	personService foundationservices.PersonService
}

func NewPersonClient(personService foundationservices.PersonService) foundationclient.Client {
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
	request coresearch.Request,
) (coresearch.Result[foundationresponse.Person], error) {
	return client.personService.Search(ctx, request)
}

func (client *personClient) DeactivatePerson(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	return client.personService.Deactivate(ctx, request)
}
