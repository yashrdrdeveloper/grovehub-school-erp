package impl

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"

 coretransaction "grovehub/libs/core/database/transaction"

	"grovehub/services/identity/entity"
	accountevents "grovehub/services/identity/events/accounts"
	eventhelpers "grovehub/services/identity/events/accounts/helpers"
	identityservices "grovehub/services/identity/services"
	infrastructureclient "grovehub/services/infrastructure/client"
)

type AccountService struct {
	delegate        identityservices.AccountService
	messagingClient infrastructureclient.MessagingClient
	logger          *slog.Logger
}

func NewAccountService(
	delegate identityservices.AccountService,
	messagingClient infrastructureclient.MessagingClient,
	logger *slog.Logger,
) identityservices.AccountService {
	return &AccountService{
		delegate:        delegate,
		messagingClient: messagingClient,
		logger:          logger,
	}
}

func (service *AccountService) Create(
	ctx context.Context,
	input identityservices.CreateAccountInput,
) (entity.Account, bool, error) {
	account, created, err := service.delegate.Create(ctx, input)
	if err != nil || !created {
		return account, created, err
	}

	service.publishCreated(ctx, account)
	return account, true, nil
}

func (service *AccountService) GetById(ctx context.Context, accountId string) (entity.Account, error) {
	return service.delegate.GetById(ctx, accountId)
}

func (service *AccountService) GetByPersonGid(ctx context.Context, personGid string) (entity.Account, error) {
	return service.delegate.GetByPersonGid(ctx, personGid)
}

func (service *AccountService) GetByEmail(ctx context.Context, email string) (entity.Account, error) {
	return service.delegate.GetByEmail(ctx, email)
}

func (service *AccountService) GetByRegistrationIdempotencyKey(
	ctx context.Context,
	key string,
) (entity.Account, error) {
	return service.delegate.GetByRegistrationIdempotencyKey(ctx, key)
}

func (service *AccountService) Update(
	ctx context.Context,
	accountId string,
	input identityservices.UpdateAccountInput,
) (entity.Account, error) {
	return service.delegate.Update(ctx, accountId, input)
}

func (service *AccountService) Deactivate(ctx context.Context, accountId string) (entity.Account, error) {
	return service.delegate.Deactivate(ctx, accountId)
}

func (service *AccountService) PermanentlyDelete(ctx context.Context, accountId string) error {
	return service.delegate.PermanentlyDelete(ctx, accountId)
}

func (service *AccountService) publishCreated(ctx context.Context, account entity.Account) {
	if coretransaction.AfterCommit(ctx, func() {
        service.publishCreated(context.WithoutCancel(ctx), account)
    }) { return }
	eventID := uuid.NewString()
	envelope := eventhelpers.BuildCreated(eventID, account)
	payload, err := json.Marshal(envelope)
	if err != nil {
		service.logger.ErrorContext(ctx, "failed to marshal identity account created event", "error", err)
		return
	}

	if err := service.messagingClient.Publish(
		ctx,
		accountevents.CreatedTopic,
		payload,
		eventhelpers.CreatedMetadata(eventID),
	); err != nil {
		service.logger.ErrorContext(
			ctx,
			"failed to publish identity account created event",
			"eventId", eventID,
			"accountId", string(account.GetAccountId()),
			"error", err,
		)
	}
}
