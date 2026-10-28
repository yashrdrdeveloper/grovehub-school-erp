package impl

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"grovehub/services/identity/entity"
	identityrepository "grovehub/services/identity/repository"
)

const (
	passwordCredentialPrimaryConstraint   = "pk_password_credentials"
	passwordCredentialAccountFKConstraint = "fk_password_credentials_account"
)

type PasswordCredentialRepository struct {
	db *gorm.DB
}

func NewPasswordCredentialRepository(db *gorm.DB) identityrepository.PasswordCredentialRepository {
	return &PasswordCredentialRepository{db: db}
}

func (repository *PasswordCredentialRepository) Create(
	ctx context.Context,
	credential entity.PasswordCredential,
) error {
	return mapPasswordCredentialWriteError(repository.db.WithContext(ctx).Create(&credential).Error)
}

func (repository *PasswordCredentialRepository) FindByAccountId(
	ctx context.Context,
	accountId entity.AccountId,
) (entity.PasswordCredential, error) {
	var credential entity.PasswordCredential
	if err := repository.db.WithContext(ctx).
		Where(&entity.PasswordCredential{AccountId: accountId}).
		Take(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.PasswordCredential{}, identityrepository.ErrPasswordCredentialNotFound
		}
		return entity.PasswordCredential{}, err
	}
	return credential, nil
}

func (repository *PasswordCredentialRepository) Update(
	ctx context.Context,
	credential entity.PasswordCredential,
) error {
	expectedVersion := credential.GetVersion() - 1
	result := repository.db.WithContext(ctx).
		Model(&entity.PasswordCredential{}).
		Where(&entity.PasswordCredential{AccountId: credential.GetAccountId()}).
		Where("version = ?", expectedVersion).
		Select("PasswordHash", "PasswordChangedAt", "Version", "UpdatedAt").
		Updates(&credential)
	if result.Error != nil {
		return mapPasswordCredentialWriteError(result.Error)
	}
	if result.RowsAffected == 0 {
		return identityrepository.ErrPasswordCredentialVersionConflict
	}
	return nil
}

func mapPasswordCredentialWriteError(err error) error {
	if err == nil {
		return nil
	}

	var postgresError *pgconn.PgError
	if !errors.As(err, &postgresError) {
		return err
	}

	if postgresError.Code == "23505" && postgresError.ConstraintName == passwordCredentialPrimaryConstraint {
		return identityrepository.ErrPasswordCredentialAlreadyExists
	}
	if postgresError.Code == "23503" && postgresError.ConstraintName == passwordCredentialAccountFKConstraint {
		return identityrepository.ErrPasswordCredentialAccountNotFound
	}
	return err
}
