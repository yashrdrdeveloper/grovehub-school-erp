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
	servicehelpers "grovehub/services/identity/services/helpers"
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
	return &AccountService{repository: repository, validator: validator}
}

func (service *AccountService) Create(
	ctx context.Context,
	input identityservices.CreateAccountInput,
) (entity.Account, bool, error) {
	validatedInput, err := service.validator.ValidateCreate(input)
	if err != nil {
		return entity.Account{}, false, err
	}

	existing, err := service.repository.FindByPersonGid(ctx, entity.PersonGid(validatedInput.PersonGid))
	if err == nil {
		if accountMatchesCreateInput(existing, validatedInput) {
			return existing, false, nil
		}
		return entity.Account{}, false, coreerrors.NewError(identityerror.AccountAlreadyExists)
	}
	if !errors.Is(err, identityrepository.ErrAccountNotFound) {
		return entity.Account{}, false, coreerrors.Wrap(identityerror.AccountReadFailed, err)
	}

	accountIdValue, err := coreid.NewUUIDv7()
	if err != nil {
		return entity.Account{}, false, coreerrors.Wrap(identityerror.AccountCreateFailed, err)
	}

	now := time.Now().UTC()
	account := servicehelpers.NewAccount(entity.AccountId(accountIdValue), validatedInput, now)

	if err := service.repository.Create(ctx, account); err != nil {
		if errors.Is(err, identityrepository.ErrAccountPersonAlreadyExists) {
			existing, readErr := service.repository.FindByPersonGid(ctx, account.PersonGid)
			if readErr == nil && accountMatchesCreateInput(existing, validatedInput) {
				return existing, false, nil
			}
			if readErr != nil && !errors.Is(readErr, identityrepository.ErrAccountNotFound) {
				return entity.Account{}, false, coreerrors.Wrap(identityerror.AccountReadFailed, readErr)
			}
		}
		if errors.Is(err, identityrepository.ErrAccountRegistrationIdempotencyKeyAlreadyExists) &&
			validatedInput.RegistrationIdempotencyKey != "" {
			existing, readErr := service.repository.FindByRegistrationIdempotencyKey(ctx, validatedInput.RegistrationIdempotencyKey)
			if readErr == nil {
				if accountMatchesCreateInput(existing, validatedInput) {
					return existing, false, nil
				}
				return entity.Account{}, false, coreerrors.NewError(identityerror.AccountRegistrationIdempotencyConflict)
			}
			if !errors.Is(readErr, identityrepository.ErrAccountNotFound) {
				return entity.Account{}, false, coreerrors.Wrap(identityerror.AccountReadFailed, readErr)
			}
		}
		return entity.Account{}, false, mapAccountWriteError(err, identityerror.AccountCreateFailed)
	}
	return account, true, nil
}

func (service *AccountService) GetById(ctx context.Context, accountId string) (entity.Account, error) {
	validatedAccountId, err := service.validator.ValidateAccountId(accountId)
	if err != nil {
		return entity.Account{}, err
	}
	return service.readAccount(service.repository.FindById(ctx, validatedAccountId))
}

func (service *AccountService) GetByPersonGid(ctx context.Context, personGid string) (entity.Account, error) {
	validatedPersonGid, err := service.validator.ValidatePersonGid(personGid)
	if err != nil {
		return entity.Account{}, err
	}
	return service.readAccount(service.repository.FindByPersonGid(ctx, validatedPersonGid))
}

func (service *AccountService) GetByEmail(ctx context.Context, email string) (entity.Account, error) {
	validatedEmail, err := service.validator.ValidateEmail(email)
	if err != nil {
		return entity.Account{}, err
	}
	return service.readAccount(service.repository.FindByEmail(ctx, validatedEmail))
}

func (service *AccountService) GetByRegistrationIdempotencyKey(
	ctx context.Context,
	key string,
) (entity.Account, error) {
	return service.readAccount(service.repository.FindByRegistrationIdempotencyKey(ctx, key))
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
	updatedAccount.Version = account.Version + 1
	updatedAccount.UpdatedAt = time.Now().UTC()
	if err := service.repository.Update(ctx, updatedAccount); err != nil {
		if errors.Is(err, identityrepository.ErrAccountVersionConflict) {
			latest, readErr := service.GetById(ctx, accountId)
			if readErr != nil {
				return entity.Account{}, readErr
			}
			_, stillChanged := applyAccountUpdate(latest, validatedInput)
			if !stillChanged {
				return latest, nil
			}
		}
		return entity.Account{}, mapAccountWriteError(err, identityerror.AccountUpdateFailed)
	}
	return updatedAccount, nil
}

func (service *AccountService) Deactivate(ctx context.Context, accountId string) (entity.Account, error) {
	account, err := service.GetById(ctx, accountId)
	if err != nil {
		return entity.Account{}, err
	}
	if account.Status == entity.AccountStatusDeactivated {
		return account, nil
	}
	updatedAccount := account
	updatedAccount.Status = entity.AccountStatusDeactivated
	updatedAccount.Version = account.Version + 1
	updatedAccount.UpdatedAt = time.Now().UTC()
	if err := service.repository.Update(ctx, updatedAccount); err != nil {
		if errors.Is(err, identityrepository.ErrAccountVersionConflict) {
			latest, readErr := service.GetById(ctx, accountId)
			if readErr != nil {
				return entity.Account{}, readErr
			}
			if latest.Status == entity.AccountStatusDeactivated {
				return latest, nil
			}
		}
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
			return nil
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

func applyAccountUpdate(account entity.Account, input identityservices.UpdateAccountInput) (entity.Account, bool) {
	updatedAccount := account
	changed := false
	if input.Email != nil && *input.Email != account.Email {
		updatedAccount.Email = *input.Email
		changed = true
	}
	if input.Phone != nil && *input.Phone != servicehelpers.NullableString(account.Phone) {
		updatedAccount.Phone = servicehelpers.OptionalString(*input.Phone)
		changed = true
	}
	return updatedAccount, changed
}

func accountMatchesCreateInput(account entity.Account, input identityservices.CreateAccountInput) bool {
	return string(account.PersonGid) == input.PersonGid &&
		account.Email == input.Email &&
		servicehelpers.NullableString(account.Phone) == input.Phone &&
		servicehelpers.NullableString(account.RegistrationIdempotencyKey) == input.RegistrationIdempotencyKey &&
		servicehelpers.NullableString(account.RegistrationRequestHash) == input.RegistrationRequestHash &&
		account.Status == entity.AccountStatusActive
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
		{identityrepository.ErrAccountRegistrationIdempotencyKeyAlreadyExists, identityerror.AccountAlreadyExists},
		{identityrepository.ErrAccountVersionConflict, identityerror.AccountConcurrentModification},
	}
	for _, item := range definitions {
		if errors.Is(err, item.target) {
			return coreerrors.NewError(item.definition)
		}
	}
	return coreerrors.Wrap(fallback, err)
}
