package impl

import (
	"context"
	"errors"
	"time"

	coreerrors "grovehub/libs/core/errors"
	identityerror "grovehub/services/identity/constants/error"
	"grovehub/services/identity/entity"
	identityrepository "grovehub/services/identity/repository"
	identityservices "grovehub/services/identity/services"
	servicehelpers "grovehub/services/identity/services/helpers"
	identityvalidators "grovehub/services/identity/validators"
)

type PasswordCredentialService struct {
	repository identityrepository.PasswordCredentialRepository
	validator  identityvalidators.PasswordCredentialValidator
}

func NewPasswordCredentialService(
	repository identityrepository.PasswordCredentialRepository,
	validator identityvalidators.PasswordCredentialValidator,
) identityservices.PasswordCredentialService {
	return &PasswordCredentialService{
		repository: repository,
		validator:  validator,
	}
}

func (service *PasswordCredentialService) Set(
	ctx context.Context,
	accountId string,
	password string,
) error {
	validatedAccountId, err := service.validatePasswordInput(accountId, password)
	if err != nil {
		return err
	}

	passwordHash, err := servicehelpers.HashPassword(password)
	if err != nil {
		return coreerrors.Wrap(identityerror.PasswordCredentialCreateFailed, err)
	}

	now := time.Now().UTC()
	credential := entity.NewPasswordCredential(
		validatedAccountId,
		passwordHash,
		now,
		1,
		now,
		now,
	)
	if err := service.repository.Create(ctx, credential); err != nil {
		return mapPasswordCredentialWriteError(err, identityerror.PasswordCredentialCreateFailed)
	}
	return nil
}

func (service *PasswordCredentialService) Verify(
	ctx context.Context,
	accountId string,
	password string,
) (bool, error) {
	validatedAccountId, err := service.validatePasswordCandidate(accountId, password)
	if err != nil {
		return false, err
	}

	credential, err := service.readCredential(ctx, validatedAccountId)
	if err != nil {
		return false, err
	}

	verified, err := servicehelpers.VerifyPassword(password, credential.GetPasswordHash())
	if err != nil {
		return false, coreerrors.Wrap(identityerror.PasswordCredentialVerifyFailed, err)
	}
	return verified, nil
}

func (service *PasswordCredentialService) Replace(
	ctx context.Context,
	accountId string,
	password string,
) error {
	validatedAccountId, err := service.validatePasswordInput(accountId, password)
	if err != nil {
		return err
	}

	credential, err := service.readCredential(ctx, validatedAccountId)
	if err != nil {
		return err
	}

	passwordHash, err := servicehelpers.HashPassword(password)
	if err != nil {
		return coreerrors.Wrap(identityerror.PasswordCredentialUpdateFailed, err)
	}

	now := time.Now().UTC()
	updatedCredential := credential.
		SetPasswordHash(passwordHash).
		SetPasswordChangedAt(now).
		SetVersion(credential.GetVersion() + 1).
		SetUpdatedAt(now)

	if err := service.repository.Update(ctx, updatedCredential); err != nil {
		return mapPasswordCredentialWriteError(err, identityerror.PasswordCredentialUpdateFailed)
	}
	return nil
}

func (service *PasswordCredentialService) validatePasswordInput(
	accountId string,
	password string,
) (entity.AccountId, error) {
	validatedAccountId, err := service.validator.ValidateAccountId(accountId)
	if err != nil {
		return "", err
	}
	if err := service.validator.ValidatePassword(password); err != nil {
		return "", err
	}
	return validatedAccountId, nil
}

func (service *PasswordCredentialService) validatePasswordCandidate(
	accountId string,
	password string,
) (entity.AccountId, error) {
	validatedAccountId, err := service.validator.ValidateAccountId(accountId)
	if err != nil {
		return "", err
	}
	if err := service.validator.ValidatePasswordCandidate(password); err != nil {
		return "", err
	}
	return validatedAccountId, nil
}

func (service *PasswordCredentialService) readCredential(
	ctx context.Context,
	accountId entity.AccountId,
) (entity.PasswordCredential, error) {
	credential, err := service.repository.FindByAccountId(ctx, accountId)
	if errors.Is(err, identityrepository.ErrPasswordCredentialNotFound) {
		return entity.PasswordCredential{}, coreerrors.NewError(identityerror.PasswordCredentialNotFound)
	}
	if err != nil {
		return entity.PasswordCredential{}, coreerrors.Wrap(identityerror.PasswordCredentialReadFailed, err)
	}
	return credential, nil
}

func mapPasswordCredentialWriteError(err error, fallback coreerrors.Definition) error {
	mapped := map[error]coreerrors.Definition{
		identityrepository.ErrPasswordCredentialAlreadyExists:  identityerror.PasswordCredentialAlreadyExists,
		identityrepository.ErrPasswordCredentialNotFound:       identityerror.PasswordCredentialNotFound,
		identityrepository.ErrPasswordCredentialAccountNotFound: identityerror.AccountNotFound,
	}
	for repositoryError, definition := range mapped {
		if errors.Is(err, repositoryError) {
			return coreerrors.NewError(definition)
		}
	}
	return coreerrors.Wrap(fallback, err)
}
