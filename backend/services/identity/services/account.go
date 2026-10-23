package services

import (
	"context"
	"errors"
	"time"

	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/helpers/id"
	foundationclient "grovehub/services/foundation/client"
	foundationerror "grovehub/services/foundation/constants/error"
	identityerror "grovehub/services/identity/constants/error"
	identityrequest "grovehub/services/identity/dto/request"
	identityresponse "grovehub/services/identity/dto/response"
	"grovehub/services/identity/entity"
	"grovehub/services/identity/mappers"
	identityrepository "grovehub/services/identity/repository"
	identitysecurity "grovehub/services/identity/security"
	servicehelpers "grovehub/services/identity/services/helpers"
)

type AccountService struct {
	repository       identityrepository.AccountRepository
	foundationClient foundationclient.Client
	passwordHasher   identitysecurity.PasswordHasher
}

func NewAccountService(
	repository identityrepository.AccountRepository,
	foundationClient foundationclient.Client,
	passwordHasher identitysecurity.PasswordHasher,
) *AccountService {
	return &AccountService{
		repository:       repository,
		foundationClient: foundationClient,
		passwordHasher:   passwordHasher,
	}
}

func (service *AccountService) Bootstrap(
	ctx context.Context,
	request identityrequest.BootstrapAccount,
) (identityresponse.BootstrapAccount, error) {
	identifierType, normalizedIdentifier, err := servicehelpers.ValidateBootstrapAccount(request)
	if err != nil {
		return identityresponse.BootstrapAccount{}, err
	}

	hasAccounts, err := service.repository.HasAccounts(ctx)
	if err != nil {
		return identityresponse.BootstrapAccount{}, coreerrors.Wrap(identityerror.BootstrapCreateFailed, err)
	}
	if hasAccounts {
		return identityresponse.BootstrapAccount{}, coreerrors.NewError(identityerror.BootstrapAlreadyCompleted)
	}

	secretHash, err := service.passwordHasher.Hash(request.Password)
	if err != nil {
		return identityresponse.BootstrapAccount{}, coreerrors.Wrap(identityerror.BootstrapCreateFailed, err)
	}

	accountId, identifierId, credentialId, err := newBootstrapIds()
	if err != nil {
		return identityresponse.BootstrapAccount{}, coreerrors.Wrap(identityerror.BootstrapCreateFailed, err)
	}

	person, _, err := service.foundationClient.CreatePerson(ctx, request.Person)
	if err != nil {
		return identityresponse.BootstrapAccount{}, mapFoundationBootstrapError(err)
	}

	now := time.Now().UTC()
	account := entity.NewAccount(
		entity.AccountId(accountId),
		entity.PersonGid(person.GID),
		entity.AccountStatusActive,
		1,
		now,
		now,
	)
	identifier := entity.NewLoginIdentifier(
		entity.LoginIdentifierId(identifierId),
		account.GetAccountId(),
		identifierType,
		normalizedIdentifier,
		normalizedIdentifier,
		true,
		servicehelpers.PersonIdentifierVerified(person, identifierType, normalizedIdentifier),
		entity.LoginIdentifierStatusActive,
		now,
		now,
	)
	credential := entity.NewCredential(
		entity.CredentialId(credentialId),
		account.GetAccountId(),
		entity.CredentialTypePassword,
		secretHash,
		entity.CredentialAlgorithmArgon2Id,
		entity.CredentialStatusActive,
		nil,
		now,
		now,
	)

	if err := service.repository.CreateFirst(
		ctx,
		account,
		[]entity.LoginIdentifier{identifier},
		[]entity.Credential{credential},
	); err != nil {
		if errors.Is(err, identityrepository.ErrBootstrapAlreadyCompleted) {
			return identityresponse.BootstrapAccount{}, coreerrors.NewError(identityerror.BootstrapAlreadyCompleted)
		}
		return identityresponse.BootstrapAccount{}, coreerrors.Wrap(identityerror.BootstrapCreateFailed, err)
	}

	return mappers.ToBootstrapAccountResponse(account, identifier), nil
}

func newBootstrapIds() (string, string, string, error) {
	values := make([]string, 3)
	for index := range values {
		value, err := coreid.NewUUIDv7()
		if err != nil {
			return "", "", "", err
		}
		values[index] = value
	}
	return values[0], values[1], values[2], nil
}

func mapFoundationBootstrapError(err error) error {
	var applicationError *coreerrors.Error
	if errors.As(err, &applicationError) {
		invalidDefinitions := map[coreerrors.Code]struct{}{
			foundationerror.InvalidPersonRequest.Code():  {},
			foundationerror.InvalidPersonContact.Code():  {},
			foundationerror.PersonContactConflict.Code(): {},
		}
		if _, ok := invalidDefinitions[applicationError.Definition().Code()]; ok {
			return coreerrors.Wrap(identityerror.InvalidBootstrapRequest, err)
		}
	}
	return coreerrors.Wrap(identityerror.BootstrapCreateFailed, err)
}
