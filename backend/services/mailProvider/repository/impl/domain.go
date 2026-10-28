package impl

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"grovehub/services/mailProvider/entity"
	mailproviderrepository "grovehub/services/mailProvider/repository"
)

const uniqueDomainNameConstraint = "uq_mail_provider_domains_name"

type DomainRepository struct {
	db *gorm.DB
}

func NewDomainRepository(db *gorm.DB) mailproviderrepository.DomainRepository {
	return &DomainRepository{db: db}
}

func (repository *DomainRepository) Create(ctx context.Context, domain entity.Domain) error {
	return mapDomainWriteError(repository.db.WithContext(ctx).Create(&domain).Error)
}

func (repository *DomainRepository) FindById(ctx context.Context, domainId entity.DomainId) (entity.Domain, error) {
	var domain entity.Domain
	if err := repository.db.WithContext(ctx).
		Where(&entity.Domain{DomainId: domainId}).
		Take(&domain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Domain{}, mailproviderrepository.ErrDomainNotFound
		}
		return entity.Domain{}, err
	}
	return domain, nil
}

func (repository *DomainRepository) FindByName(ctx context.Context, name string) (entity.Domain, error) {
	var domain entity.Domain
	if err := repository.db.WithContext(ctx).
		Where(&entity.Domain{Name: name}).
		Take(&domain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Domain{}, mailproviderrepository.ErrDomainNotFound
		}
		return entity.Domain{}, err
	}
	return domain, nil
}

func (repository *DomainRepository) Update(ctx context.Context, domain entity.Domain) error {
	result := repository.db.WithContext(ctx).
		Model(&entity.Domain{DomainId: domain.GetDomainId()}).
		Select("Status", "Version", "UpdatedAt").
		Updates(&domain)
	if result.Error != nil {
		return mapDomainWriteError(result.Error)
	}
	if result.RowsAffected == 0 {
		return mailproviderrepository.ErrDomainNotFound
	}
	return nil
}

func mapDomainWriteError(err error) error {
	if err == nil {
		return nil
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" && postgresError.ConstraintName == uniqueDomainNameConstraint {
		return mailproviderrepository.ErrDomainAlreadyExists
	}
	return err
}
