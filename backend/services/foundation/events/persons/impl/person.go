package impl

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"

	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	personevents "grovehub/services/foundation/events/persons"
	eventhelpers "grovehub/services/foundation/events/persons/helpers"
	foundationservices "grovehub/services/foundation/services"
	infrastructureclient "grovehub/services/infrastructure/client"
)

type PersonService struct {
	delegate        foundationservices.PersonService
	messagingClient infrastructureclient.MessagingClient
	logger          *slog.Logger
}

func NewPersonService(
	delegate foundationservices.PersonService,
	messagingClient infrastructureclient.MessagingClient,
	logger *slog.Logger,
) foundationservices.PersonService {
	return &PersonService{
		delegate:        delegate,
		messagingClient: messagingClient,
		logger:          logger,
	}
}

func (service *PersonService) Create(
	ctx context.Context,
	request foundationrequest.CreatePerson,
) (foundationresponse.Person, bool, error) {
	person, created, err := service.delegate.Create(ctx, request)
	if err != nil || !created {
		return person, created, err
	}

	service.publishCreated(ctx, person)
	return person, true, nil
}

func (service *PersonService) GetByGID(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	return service.delegate.GetByGID(ctx, request)
}

func (service *PersonService) Update(
	ctx context.Context,
	personRequest foundationrequest.PersonByGID,
	request foundationrequest.UpdatePerson,
) (foundationresponse.Person, error) {
	return service.delegate.Update(ctx, personRequest, request)
}

func (service *PersonService) Search(
	ctx context.Context,
	request coresearch.Request,
) (coresearch.Result[foundationresponse.Person], error) {
	return service.delegate.Search(ctx, request)
}

func (service *PersonService) Deactivate(
	ctx context.Context,
	request foundationrequest.PersonByGID,
) (foundationresponse.Person, error) {
	return service.delegate.Deactivate(ctx, request)
}

func (service *PersonService) publishCreated(ctx context.Context, person foundationresponse.Person) {
	eventID := uuid.NewString()
	envelope := eventhelpers.BuildCreated(eventID, person)
	payload, err := json.Marshal(envelope)
	if err != nil {
		service.logger.ErrorContext(ctx, "failed to marshal foundation person created event", "error", err)
		return
	}

	if err := service.messagingClient.Publish(
		ctx,
		personevents.CreatedTopic,
		payload,
		eventhelpers.CreatedMetadata(eventID),
	); err != nil {
		service.logger.ErrorContext(
			ctx,
			"failed to publish foundation person created event",
			"eventId", eventID,
			"gid", person.GID,
			"error", err,
		)
	}
}
