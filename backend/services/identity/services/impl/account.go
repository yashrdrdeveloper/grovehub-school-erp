package impl

import (
	"context"
	"errors"
	"time"

	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/id"
	identityerror "grovehub/services/identity/constants/error"
	"grovehub/services/identity/entity"
	identityrepository "grovehub/services/identity/repository"
	identityservices "grovehub/services/identity/services"
	identityvalidators "grovehub/services/identity/validators"
)

type AccountService struct {
	repository identityrepository.AccountRepository
	validator  identityvalidators.AccountValidator
}

func NewAccountService(
	repository identityrepository.AccountRepository,
	validator identityvalidators.AccountValidator,
) identityservices.AccountService {
	return &AccountService{
		repository: repository,
		validator:  validator,
	}
}

func (service *AccountService) Create(
	ctx context.Context,
	input identityservices.CreateAccountInput,
) (entity.Account, error) {
	validatedInput, err := service.validator.ValidateCreate(input)
	if err != nil {
		return entity.Account{}, err
	}

	accountIdValue, err := coreid.NewUUIDv7()
	if err != nil {
		return entity.Account{}, coreerrors.Wrap(identityerror.AccountCreateFailed, err)
	}

	now := time.Now().UTC()
	account := entity.NewAccount(
		entity.AccountId(accountIdValue),
		entity.PersonGid(validatedInput.PersonGid),
		validatedInput.Email,
		validatedInput.Phone,
		entity.AccountStatusActive,
		1,
		now,
		now,
	)

	if err := service.repository.Create(ctx, account); err != nil {
		return entity.Account{}, mapAccountWriteError(err, identityerror.AccountCreateFailed)
	}
	return account, nil
}

func (service *AccountService) GetById(
	ctx context.Context,
	accountId string,
) (entity.Account, error) {
	validatedAccountId, err := service.validator.ValidateAccountId(accountId)
	if err != nil {
		return entity.Account{}, err
	}
	return service.readAccount(service.repository.FindById(ctx, validatedAccountId))
}

func (service *AccountService) GetByPersonGid(
	ctx context.Context,
	personGid string,
) (entity.Account, error) {
	validatedPersonGid, err := service.validator.ValidatePersonGid(personGid)
	if err != nil {
		return entity.Account{}, err
	}
	return service.readAccount(service.repository.FindByPersonGid(ctx, validatedPersonGid))
}

func (service *AccountService) GetByEmail(
	ctx context.Context,
	email string,
) (entity.Account, error) {
	validatedEmail, err := service.validator.ValidateEmail(email)
	if err != nil {
		return entity.Account{}, err
	}
	return service.readAccount(service.repository.FindByEmail(ctx, validatedEmail))
}

func (service *AccountService) Update(
	ctx context.Context,
	accountId string,
	input identityservices.UpdateAccountInput,
) (entity.Account, error) {
	account, err := service.GetById(ctx, accountId)
	if err != nil {
		return entity.Account{}, err
	}

	validatedInput, err := service.validator.ValidateUpdate(input)
	if err != nil {
		return entity.Account{}, err
	}

	updatedAccount, changed := applyAccountUpdate(account, validatedInput)
	if !changed {
		return account, nil
	}

	updatedAccount = updatedAccount.
		SetVersion(account.GetVersion() + 1).
		SetUpdatedAt(time.Now().UTC())
	if err := service.repository.Update(ctx, updatedAccount); err != nil {
		return entity.Account{}, mapAccountWriteError(err, identityerror.AccountUpdateFailed)
	}
	return updatedAccount, nil
}

func (service *AccountService) Deactivate(
	ctx context.Context,
	accountId string,
) (entity.Account, error) {
	account, err := service.GetById(ctx, accountId)
	if err != nil {
		return entity.Account{}, err
	}
	if account.GetStatus() == entity.AccountStatusDeactivated {
		return account, nil
	}

	updatedAccount := account.
		SetStatus(entity.AccountStatusDeactivated).
		SetVersion(account.GetVersion() + 1).
		SetUpdatedAt(time.Now().UTC())
	if err := service.repository.Update(ctx, updatedAccount); err != nil {
		return entity.Account{}, mapAccountWriteError(err, identityerror.AccountUpdateFailed)
	}
	return updatedAccount, nil
}

func (service *AccountService) PermanentlyDelete(ctx context.Context, accountId string) error {
	validatedAccountId, err := service.validator.ValidateAccountId(accountId)
	if err != nil {
		return err
	}
	if err := service.repository.Delete(ctx, validatedAccountId); err != nil {
		if errors.Is(err, identityrepository.ErrAccountNotFound) {
			return coreerrors.NewError(identityerror.AccountNotFound)
		}
		return coreerrors.Wrap(identityerror.AccountDeleteFailed, err)
	}
	return nil
}

func (service *AccountService) readAccount(account entity.Account, err error) (entity.Account, error) {
	if errors.Is(err, identityrepository.ErrAccountNotFound) {
		return entity.Account{}, coreerrors.NewError(identityerror.AccountNotFound)
	}
	if err != nil {
		return entity.Account{}, coreerrors.Wrap(identityerror.AccountReadFailed, err)
	}
	return account, nil
}

func applyAccountUpdate(
	account entity.Account,
	input identityservices.UpdateAccountInput,
) (entity.Account, bool) {
	updatedAccount := account
	changed := false

	if input.Email != nil && *input.Email != account.GetEmail() {
		updatedAccount = updatedAccount.SetEmail(*input.Email)
		changed = true
	}

	if input.Phone != nil && *input.Phone != account.GetPhone() {
		updatedAccount = updatedAccount.SetPhone(*input.Phone)
		changed = true
	}

	return updatedAccount, changed
}

func mapAccountWriteError(err error, fallback coreerrors.Definition) error {
	definitions := []struct {
		target     error
		definition coreerrors.Definition
	}{
		{identityrepository.ErrAccountNotFound, identityerror.AccountNotFound},
		{identityrepository.ErrAccountPersonAlreadyExists, identityerror.AccountAlreadyExists},
		{identityrepository.ErrAccountEmailAlreadyExists, identityerror.AccountEmailAlreadyExists},
		{identityrepository.ErrAccountPhoneAlreadyExists, identityerror.AccountPhoneAlreadyExists},
	}
	for _, item := range definitions {
		if errors.Is(err, item.target) {
			return coreerrors.NewError(item.definition)
		}
	}
	return coreerrors.Wrap(fallback, err)
}
