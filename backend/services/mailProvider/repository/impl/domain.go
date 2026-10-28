package impl

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"grovehub/services/mailProvider/entity"
	mailproviderrepository "grovehub/services/mailProvider/repository"
)

const uniqueDomainNameConstraint = "uq_mail_provider_domains_name"

type DomainRepository struct {
	db *gorm.DB
}

func NewDomainRepository(db *gorm.DB) mailproviderrepository.DomainRepository { return &DomainRepository{db: db} }

func (repository *DomainRepository) Create(ctx context.Context, domain entity.Domain) error {
	return mapDomainWriteError(repository.db.WithContext(ctx).Create(&domain).Error)
}

func (repository *DomainRepository) FindAll(ctx context.Context) ([]entity.Domain, error) {
	domains := make([]entity.Domain, 0)
	if err := repository.db.WithContext(ctx).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		Find(&domains).Error; err != nil {
		return nil, err
	}
	return domains, nil
}

func (repository *DomainRepository) FindById(ctx context.Context, domainId entity.DomainId) (entity.Domain, error) {
	var domain entity.Domain
	if err := repository.db.WithContext(ctx).Where(&entity.Domain{DomainId: domainId}).Take(&domain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Domain{}, mailproviderrepository.ErrDomainNotFound
		}
		return entity.Domain{}, err
	}
	return domain, nil
}

func (repository *DomainRepository) FindByName(ctx context.Context, name string) (entity.Domain, error) {
	var domain entity.Domain
	if err := repository.db.WithContext(ctx).Where(&entity.Domain{Name: name}).Take(&domain).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Domain{}, mailproviderrepository.ErrDomainNotFound
		}
		return entity.Domain{}, err
	}
	return domain, nil
}

func (repository *DomainRepository) Update(ctx context.Context, domain entity.Domain) error {
	expectedVersion := domain.GetVersion() - 1
	result := repository.db.WithContext(ctx).
		Model(&entity.Domain{}).
		Where(&entity.Domain{DomainId: domain.GetDomainId()}).
		Where("version = ?", expectedVersion).
		Select("Status", "Version", "UpdatedAt").
		Updates(&domain)
	if result.Error != nil {
		return mapDomainWriteError(result.Error)
	}
	if result.RowsAffected == 0 {
		return mailproviderrepository.ErrDomainVersionConflict
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
