package impl

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"grovehub/services/infrastructure/entity"
	infrastructurerepository "grovehub/services/infrastructure/repository"
)

type DKIMRepository struct {
	db *gorm.DB
}

func NewDKIMRepository(db *gorm.DB) infrastructurerepository.DKIMRepository {
	return &DKIMRepository{db: db}
}

func (repository *DKIMRepository) Create(ctx context.Context, key entity.DKIMKey) error {
	err := repository.db.WithContext(ctx).Create(&key).Error
	if err == nil {
		return nil
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		return infrastructurerepository.ErrDKIMKeyAlreadyExists
	}
	return err
}

func (repository *DKIMRepository) FindByDomain(ctx context.Context, domain string) (entity.DKIMKey, error) {
	var key entity.DKIMKey
	if err := repository.db.WithContext(ctx).Where(&entity.DKIMKey{Domain: domain}).Take(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.DKIMKey{}, infrastructurerepository.ErrDKIMKeyNotFound
		}
		return entity.DKIMKey{}, err
	}
	return key, nil
}
