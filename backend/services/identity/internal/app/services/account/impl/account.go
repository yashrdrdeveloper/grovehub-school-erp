package impl

import (
	"context"
	"errors"
	"strings"
	"time"

	transaction "grovehub/libs/core/database/transaction"
	coresearch "grovehub/libs/core/search"
	foundation "grovehub/services/foundation/client/person"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	accountproducer "grovehub/services/identity/internal/app/events/producers/account"
	accountmapper "grovehub/services/identity/internal/app/mappers/account"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	"grovehub/services/identity/internal/app/security/principal"
	account "grovehub/services/identity/internal/app/services/account"
	"grovehub/services/identity/internal/app/services/account/helpers"
	create "grovehub/services/identity/internal/app/services/account/operations/create"
	search "grovehub/services/identity/internal/app/services/account/operations/search"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	credential "grovehub/services/identity/internal/app/services/credential"
	invitation "grovehub/services/identity/internal/app/services/registration/operations/invitation"
	accountvalidator "grovehub/services/identity/internal/app/validators/account"
)

// @Service
type Service struct {
	accounts accountrepo.Repository
	tokens tokenrepo.Repository
	sessions sessionrepo.Repository
	credentials credential.Service
	invitations *invitation.Service
	mapper accountmapper.Mapper
	validator accountvalidator.Validator
	unit transaction.UnitOfWork
	events accountproducer.ChangedProducer
	create *create.Operation
	search *search.Operation
}

var _ account.Service = (*Service)(nil)

func New(accounts accountrepo.Repository, tokens tokenrepo.Repository, sessions sessionrepo.Repository, foundationClient foundation.Client, credentials credential.Service, invitations *invitation.Service, mapper accountmapper.Mapper, validator accountvalidator.Validator, unit transaction.UnitOfWork, events accountproducer.ChangedProducer) account.Service {
	return &Service{
		accounts: accounts,
		tokens: tokens,
		sessions: sessions,
		credentials: credentials,
		invitations: invitations,
		mapper: mapper,
		validator: validator,
		unit: unit,
		events: events,
		create: create.New(accounts, foundationClient, validator, invitations, unit, events),
		search: search.New(accounts, mapper, validator),
	}
}

func projection(account entity.Account) string {
	username := ""
	if account.Username != nil {
		username = *account.Username
	}
	return strings.ToLower(strings.Join([]string{account.Email, username}, " "))
}

func (service *Service) Create(ctx context.Context, input accountdto.Create, idempotencyKey string) (accountdto.CreateResult, error) {
	return service.create.Execute(ctx, input, idempotencyKey)
}

func (service *Service) Search(ctx context.Context, request coresearch.Request) (coresearch.Result[accountdto.Account], error) {
	return service.search.Execute(ctx, request)
}

func (service *Service) Get(ctx context.Context, accountID string) (accountdto.Account, error) {
	actor, err := principal.RequireAdmin(ctx)
	if err != nil {
		return accountdto.Account{}, err
	}
	accountEntity, err := service.accounts.FindByID(ctx, accountID)
	if err != nil {
		return accountdto.Account{}, err
	}
	if accountEntity.TenantID == nil || *accountEntity.TenantID != actor.TenantID {
		return accountdto.Account{}, accountrepo.ErrNotFound
	}
	return service.mapper.ToResponse(accountEntity), nil
}

func (service *Service) Update(ctx context.Context, accountID string, request accountdto.Update) (accountdto.Account, error) {
	actor, err := principal.RequireAdmin(ctx)
	if err != nil {
		return accountdto.Account{}, err
	}
	if request.Username == nil {
		return accountdto.Account{}, errors.New("username must be supplied")
	}
	username, err := service.validator.NormalizeUsername(*request.Username)
	if err != nil {
		return accountdto.Account{}, err
	}

	var updated entity.Account
	err = service.unit.Run(ctx, func(transactionCtx context.Context) error {
		current, err := service.accounts.FindByID(transactionCtx, accountID)
		if err != nil {
			return err
		}
		if current.TenantID == nil || *current.TenantID != actor.TenantID {
			return accountrepo.ErrNotFound
		}

		before := current
		current.Username = &username
		current.Version++
		current.UpdatedAt = time.Now().UTC()
		current.SearchDetails = projection(current)
		if err := service.accounts.Update(transactionCtx, current); err != nil {
			return err
		}
		if err := service.events.Enqueue(transactionCtx, "UPDATED", &before, &current); err != nil {
			return err
		}
		updated = current
		return nil
	})
	if err != nil {
		return accountdto.Account{}, err
	}
	return service.mapper.ToResponse(updated), nil
}

func (service *Service) Transition(ctx context.Context, accountID string, action string) (accountdto.Account, error) {
	actor, err := principal.RequireAdmin(ctx)
	if err != nil {
		return accountdto.Account{}, err
	}

	var updated entity.Account
	err = service.unit.Run(ctx, func(transactionCtx context.Context) error {
		current, err := service.accounts.FindByID(transactionCtx, accountID)
		if err != nil {
			return err
		}
		if current.TenantID == nil || *current.TenantID != actor.TenantID {
			return accountrepo.ErrNotFound
		}

		next, err := helpers.NextStatus(current.Status, action)
		if err != nil {
			return err
		}

		before := current
		current.Status = next
		current.Version++
		current.UpdatedAt = time.Now().UTC()
		if err := service.accounts.Update(transactionCtx, current); err != nil {
			return err
		}

		if next == entity.Suspended || next == entity.Deactivated {
			if _, err := service.sessions.RevokeAll(transactionCtx, current.AccountID); err != nil {
				return err
			}
		}

		operation := map[string]string{
			"suspend": "SUSPENDED",
			"activate": "ACTIVATED",
			"deactivate": "DEACTIVATED",
		}[action]
		if err := service.events.Enqueue(transactionCtx, operation, &before, &current); err != nil {
			return err
		}
		updated = current
		return nil
	})
	if err != nil {
		return accountdto.Account{}, err
	}
	return service.mapper.ToResponse(updated), nil
}

func (service *Service) Resend(ctx context.Context, accountID string) error {
	actor, err := principal.RequireAdmin(ctx)
	if err != nil {
		return err
	}

	return service.unit.Run(ctx, func(transactionCtx context.Context) error {
		accountEntity, err := service.accounts.FindByID(transactionCtx, accountID)
		if err != nil {
			return err
		}
		if accountEntity.TenantID == nil || *accountEntity.TenantID != actor.TenantID {
			return accountrepo.ErrNotFound
		}
		if accountEntity.Status != entity.PendingSetup {
			return errors.New("account is not pending setup")
		}
		return service.invitations.Queue(transactionCtx, accountEntity, entity.AccountSetup, "there")
	})
}

func (service *Service) Setup(ctx context.Context, grant string, input accountdto.Setup) (accountdto.SetupResult, error) {
	username, err := service.validator.NormalizeUsername(input.Username)
	if err != nil {
		return accountdto.SetupResult{}, err
	}
	digest, err := secrets.Digest(grant, 32)
	if err != nil {
		return accountdto.SetupResult{}, err
	}

	var result accountdto.SetupResult
	err = service.unit.Run(ctx, func(transactionCtx context.Context) error {
		accountID, err := service.tokens.Consume(transactionCtx, digest, entity.SetupGrant, time.Now().UTC())
		if err != nil {
			return err
		}
		accountEntity, err := service.accounts.FindByID(transactionCtx, accountID)
		if err != nil {
			return err
		}
		if accountEntity.Status != entity.PendingSetup {
			return errors.New("account not pending setup")
		}

		before := accountEntity
		accountEntity.Username = &username
		accountEntity.Status = entity.Active
		accountEntity.Version++
		accountEntity.UpdatedAt = time.Now().UTC()
		accountEntity.SearchDetails = projection(accountEntity)

		if err := service.credentials.Set(transactionCtx, accountID, input.Password); err != nil {
			return err
		}
		if err := service.accounts.Update(transactionCtx, accountEntity); err != nil {
			return err
		}
		if err := service.events.Enqueue(transactionCtx, "ACTIVATED", &before, &accountEntity); err != nil {
			return err
		}

		result = accountdto.SetupResult{
			AccountID: accountID,
			PersonGID: accountEntity.PersonGID,
			Username: username,
			Status: string(accountEntity.Status),
			SetupCompleted: true,
		}
		return nil
	})
	return result, err
}
