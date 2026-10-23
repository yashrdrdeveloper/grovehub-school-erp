package client

import (
	"context"

	searchcontract "grovehub/libs/core/contracts/search"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
)

type Client interface {
	CreatePerson(
		ctx context.Context,
		request foundationrequest.CreatePerson,
	) (foundationresponse.Person, bool, error)
	GetPersonByGID(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
	UpdatePerson(
		ctx context.Context,
		personRequest foundationrequest.PersonByGID,
		request foundationrequest.UpdatePerson,
	) (foundationresponse.Person, error)
	SearchPersons(
		ctx context.Context,
		request searchcontract.Request,
	) (searchcontract.Result[foundationresponse.Person], error)
	DeactivatePerson(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
}

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

type foundationClient struct {
	personService personService
}

func New(personService personService) Client {
	return &foundationClient{personService: personService}
}

func (client *foundationClient) CreatePerson(
	ctx context.Context,
	request foundationrequest.CreatePerson,
) (foundationresponse.Person, bool, error) {
	return client.personService.Create(ctx, request)
}

func (client *foundationClient) GetPersonByGID(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	return client.personService.GetByGID(ctx, request)
}

func (client *foundationClient) UpdatePerson(
	ctx context.Context,
	personRequest foundationrequest.PersonByGID,
	request foundationrequest.UpdatePerson,
) (foundationresponse.Person, error) {
	return client.personService.Update(ctx, personRequest, request)
}

func (client *foundationClient) SearchPersons(
	ctx context.Context,
	request searchcontract.Request,
) (searchcontract.Result[foundationresponse.Person], error) {
	return client.personService.Search(ctx, request)
}

func (client *foundationClient) DeactivatePerson(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	return client.personService.Deactivate(ctx, request)
}
