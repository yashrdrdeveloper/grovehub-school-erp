package impl

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"

	coreevents "grovehub/libs/core/events"
	infrastructureclient "grovehub/services/infrastructure/client"
	"grovehub/services/identity/entity"
	identityevents "grovehub/services/identity/events"
	identityservices "grovehub/services/identity/services"
)

type accountEventService struct {
	delegate        identityservices.AccountService
	messagingClient infrastructureclient.MessagingClient
	logger          *slog.Logger
}

func NewAccountEventService(
	delegate identityservices.AccountService,
	messagingClient infrastructureclient.MessagingClient,
	logger *slog.Logger,
) identityservices.AccountService {
	return &accountEventService{
		delegate:        delegate,
		messagingClient: messagingClient,
		logger:          logger,
	}
}

func (service *accountEventService) Create(
	ctx context.Context,
	input identityservices.CreateAccountInput,
) (entity.Account, bool, error) {
	account, created, err := service.delegate.Create(ctx, input)
	if err != nil || !created {
		return account, created, err
	}

	service.publishAccountCreated(ctx, account)
	return account, true, nil
}

func (service *accountEventService) GetById(ctx context.Context, accountId string) (entity.Account, error) {
	return service.delegate.GetById(ctx, accountId)
}

func (service *accountEventService) GetByPersonGid(ctx context.Context, personGid string) (entity.Account, error) {
	return service.delegate.GetByPersonGid(ctx, personGid)
}

func (service *accountEventService) GetByEmail(ctx context.Context, email string) (entity.Account, error) {
	return service.delegate.GetByEmail(ctx, email)
}

func (service *accountEventService) GetByRegistrationIdempotencyKey(
	ctx context.Context,
	key string,
) (entity.Account, error) {
	return service.delegate.GetByRegistrationIdempotencyKey(ctx, key)
}

func (service *accountEventService) Update(
	ctx context.Context,
	accountId string,
	input identityservices.UpdateAccountInput,
) (entity.Account, error) {
	return service.delegate.Update(ctx, accountId, input)
}

func (service *accountEventService) Deactivate(ctx context.Context, accountId string) (entity.Account, error) {
	return service.delegate.Deactivate(ctx, accountId)
}

func (service *accountEventService) PermanentlyDelete(ctx context.Context, accountId string) error {
	return service.delegate.PermanentlyDelete(ctx, accountId)
}

func (service *accountEventService) publishAccountCreated(ctx context.Context, account entity.Account) {
	eventID := uuid.NewString()
	envelope := coreevents.Envelope[identityevents.AccountCreated]{
		EventID:    eventID,
		EventType:  identityevents.AccountCreatedType,
		Source:     identityevents.Source,
		OccurredAt: account.GetCreatedAt(),
		Data: identityevents.AccountCreated{
			AccountID: string(account.GetAccountId()),
			PersonGID: string(account.GetPersonGid()),
			Status:    string(account.GetStatus()),
			Version:   account.GetVersion(),
		},
	}

	payload, err := json.Marshal(envelope)
	if err != nil {
		service.logger.ErrorContext(ctx, "failed to marshal identity account created event", "error", err)
		return
	}

	metadata := map[string]string{
		"event-id":   eventID,
		"event-type": identityevents.AccountCreatedType,
		"source":     identityevents.Source,
	}
	if err := service.messagingClient.Publish(ctx, identityevents.AccountCreatedTopic, payload, metadata); err != nil {
		service.logger.ErrorContext(
			ctx,
			"failed to publish identity account created event",
			"eventId", eventID,
			"accountId", string(account.GetAccountId()),
			"error", err,
		)
	}
}
