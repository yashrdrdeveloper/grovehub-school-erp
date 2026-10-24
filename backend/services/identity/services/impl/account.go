package impl

import (
	"context"
	"errors"
	"time"

	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/helpers/id"
	identityerror "grovehub/services/identity/constants/error"
	"grovehub/services/identity/entity"
	identityrepository "grovehub/services/identity/repository"
	identityservices "grovehub/services/identity/services"
	servicehelpers "grovehub/services/identity/services/helpers"
)

type AccountService struct {
	repository identityrepository.AccountRepository
}

func NewAccountService(repository identityrepository.AccountRepository) identityservices.AccountService {
	return &AccountService{repository: repository}
}

func (service *AccountService) Create(
	ctx context.Context,
	input identityservices.CreateAccountInput,
) (entity.Account, error) {
	personGid, err := servicehelpers.NormalizePersonGid(input.PersonGid)
	if err != nil {
		return entity.Account{}, err
	}
	email, err := servicehelpers.NormalizeEmail(input.Email)
	if err != nil {
		return entity.Account{}, err
	}
	phone, err := servicehelpers.NormalizePhone(input.Phone)
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
		personGid,
		email,
		phone,
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
	normalizedAccountId, err := servicehelpers.NormalizeAccountId(accountId)
	if err != nil {
		return entity.Account{}, err
	}
	return service.readAccount(service.repository.FindById(ctx, normalizedAccountId))
}

func (service *AccountService) GetByPersonGid(
	ctx context.Context,
	personGid string,
) (entity.Account, error) {
	normalizedPersonGid, err := servicehelpers.NormalizePersonGid(personGid)
	if err != nil {
		return entity.Account{}, err
	}
	return service.readAccount(service.repository.FindByPersonGid(ctx, normalizedPersonGid))
}

func (service *AccountService) GetByEmail(
	ctx context.Context,
	email string,
) (entity.Account, error) {
	normalizedEmail, err := servicehelpers.NormalizeEmail(email)
	if err != nil {
		return entity.Account{}, err
	}
	return service.readAccount(service.repository.FindByEmail(ctx, normalizedEmail))
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

	updatedAccount, changed, err := applyAccountUpdate(account, input)
	if err != nil {
		return entity.Account{}, err
	}
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
	normalizedAccountId, err := servicehelpers.NormalizeAccountId(accountId)
	if err != nil {
		return err
	}
	if err := service.repository.Delete(ctx, normalizedAccountId); err != nil {
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
) (entity.Account, bool, error) {
	updatedAccount := account
	changed := false

	if input.Email != nil {
		email, err := servicehelpers.NormalizeEmail(*input.Email)
		if err != nil {
			return entity.Account{}, false, err
		}
		if email != account.GetEmail() {
			updatedAccount = updatedAccount.SetEmail(email)
			changed = true
		}
	}

	if input.Phone != nil {
		phone, err := servicehelpers.NormalizePhone(*input.Phone)
		if err != nil {
			return entity.Account{}, false, err
		}
		if phone != account.GetPhone() {
			updatedAccount = updatedAccount.SetPhone(phone)
			changed = true
		}
	}

	return updatedAccount, changed, nil
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
