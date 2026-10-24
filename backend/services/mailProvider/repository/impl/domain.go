package impl

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"grovehub/services/mailProvider/entity"
	mailproviderrepository "grovehub/services/mailProvider/repository"
)

const uniqueDomainNameConstraint = "uq_mail_provider_domains_name"

type DomainRepository struct {
	pool *pgxpool.Pool
}

func NewDomainRepository(pool *pgxpool.Pool) mailproviderrepository.DomainRepository {
	return &DomainRepository{pool: pool}
}

func (repository *DomainRepository) Create(ctx context.Context, domain entity.Domain) error {
	_, err := repository.pool.Exec(ctx, `INSERT INTO mail_provider.domains (
		domain_id, name, verification_token, status, version, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		string(domain.GetDomainId()), domain.GetName(), domain.GetVerificationToken(),
		string(domain.GetStatus()), domain.GetVersion(), domain.GetCreatedAt(), domain.GetUpdatedAt(),
	)
	return mapDomainWriteError(err)
}

func (repository *DomainRepository) FindById(ctx context.Context, domainId entity.DomainId) (entity.Domain, error) {
	return scanDomain(repository.pool.QueryRow(ctx, domainSelect+` WHERE domain_id = $1`, string(domainId)))
}

func (repository *DomainRepository) FindByName(ctx context.Context, name string) (entity.Domain, error) {
	return scanDomain(repository.pool.QueryRow(ctx, domainSelect+` WHERE name = $1`, name))
}

func (repository *DomainRepository) Update(ctx context.Context, domain entity.Domain) error {
	result, err := repository.pool.Exec(ctx, `UPDATE mail_provider.domains
	SET status = $2, version = $3, updated_at = $4 WHERE domain_id = $1`,
		string(domain.GetDomainId()), string(domain.GetStatus()), domain.GetVersion(), domain.GetUpdatedAt(),
	)
	if err != nil {
		return mapDomainWriteError(err)
	}
	if result.RowsAffected() == 0 {
		return mailproviderrepository.ErrDomainNotFound
	}
	return nil
}

const domainSelect = `SELECT domain_id, name, verification_token, status, version, created_at, updated_at FROM mail_provider.domains`

func scanDomain(row pgx.Row) (entity.Domain, error) {
	var domainId, name, token, status string
	var version int64
	var createdAt, updatedAt time.Time
	if err := row.Scan(&domainId, &name, &token, &status, &version, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Domain{}, mailproviderrepository.ErrDomainNotFound
		}
		return entity.Domain{}, err
	}
	return entity.NewDomain(entity.DomainId(domainId), name, token, entity.DomainStatus(status), version, createdAt, updatedAt), nil
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
