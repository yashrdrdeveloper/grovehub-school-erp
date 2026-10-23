package services

import (
	"context"
	"errors"
	"time"

	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/helpers/id"
	identityerror "grovehub/services/identity/constants/error"
	"grovehub/services/identity/entity"
	identityrepository "grovehub/services/identity/repository"
	servicehelpers "grovehub/services/identity/services/helpers"
)

type CreateAccountInput struct {
	PersonGid string
	Email     string
	Phone     string
}

type AccountService struct {
	repository identityrepository.AccountRepository
}

func NewAccountService(repository identityrepository.AccountRepository) *AccountService {
	return &AccountService{repository: repository}
}

func (service *AccountService) Create(
	ctx context.Context,
	input CreateAccountInput,
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
		return entity.Account{}, mapAccountCreateError(err)
	}
	return account, nil
}

func mapAccountCreateError(err error) error {
	definitions := []struct {
		target     error
		definition coreerrors.Definition
	}{
		{identityrepository.ErrAccountPersonAlreadyExists, identityerror.AccountAlreadyExists},
		{identityrepository.ErrAccountEmailAlreadyExists, identityerror.AccountEmailAlreadyExists},
		{identityrepository.ErrAccountPhoneAlreadyExists, identityerror.AccountPhoneAlreadyExists},
	}
	for _, item := range definitions {
		if errors.Is(err, item.target) {
			return coreerrors.NewError(item.definition)
		}
	}
	return coreerrors.Wrap(identityerror.AccountCreateFailed, err)
}
