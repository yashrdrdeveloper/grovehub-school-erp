package impl

import (
	"context"
	"errors"
	"time"

	"grovehub/services/identity/internal/app/entity"
	credentialrepo "grovehub/services/identity/internal/app/repository/credential"
	credential "grovehub/services/identity/internal/app/services/credential"
	"grovehub/services/identity/internal/app/services/credential/helpers"
)

type Service struct {
	repository credentialrepo.Repository
}

var _ credential.Service = (*Service)(nil)

func New(repository credentialrepo.Repository) credential.Service {
	return &Service{repository: repository}
}

func validate(password string) error {
	if len(password) < 12 || len(password) > 256 {
		return errors.New("invalid password length")
	}
	return nil
}

func (service *Service) Set(ctx context.Context, accountID string, password string) error {
	if err := validate(password); err != nil {
		return err
	}
	_, err := service.repository.Find(ctx, accountID)
	if err == nil {
		return errors.New("password already configured")
	}
	if !errors.Is(err, credentialrepo.ErrNotFound) {
		return err
	}

	hash, err := helpers.HashPassword(password)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	return service.repository.Create(ctx, entity.PasswordCredential{
		AccountID: accountID,
		PasswordHash: hash,
		Version: 1,
		CreatedAt: now,
		UpdatedAt: now,
		PasswordChangedAt: now,
	})
}

func (service *Service) Verify(ctx context.Context, accountID string, password string) (bool, error) {
	credential, err := service.repository.Find(ctx, accountID)
	if err != nil {
		return false, err
	}
	return helpers.VerifyPassword(password, credential.PasswordHash)
}

func (service *Service) Replace(ctx context.Context, accountID string, password string) error {
	if err := validate(password); err != nil {
		return err
	}
	credential, err := service.repository.Find(ctx, accountID)
	if err != nil {
		return err
	}
	hash, err := helpers.HashPassword(password)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	credential.Version++
	credential.PasswordHash = hash
	credential.PasswordChangedAt = now
	credential.UpdatedAt = now
	return service.repository.Update(ctx, credential)
}
