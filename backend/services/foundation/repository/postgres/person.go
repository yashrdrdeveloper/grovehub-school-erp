package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"grovehub/services/foundation/entity"
	foundationrepository "grovehub/services/foundation/repository"
)

type PersonRepository struct {
	pool *pgxpool.Pool
}

func NewPersonRepository(pool *pgxpool.Pool) *PersonRepository {
	return &PersonRepository{pool: pool}
}

func (repository *PersonRepository) Create(ctx context.Context, person entity.Person) error {
	_, err := repository.pool.Exec(
		ctx,
		`INSERT INTO foundation.persons (gid, status, version, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		string(person.GetGID()),
		string(person.GetStatus()),
		person.GetVersion(),
		person.GetCreatedAt(),
		person.GetUpdatedAt(),
	)

	return err
}

func (repository *PersonRepository) FindByGID(ctx context.Context, gid entity.GID) (entity.Person, error) {
	var gidValue string
	var status string
	var version int64
	var createdAt time.Time
	var updatedAt time.Time

	err := repository.pool.QueryRow(
		ctx,
		`SELECT gid, status, version, created_at, updated_at
		 FROM foundation.persons
		 WHERE gid = $1`,
		string(gid),
	).Scan(
		&gidValue,
		&status,
		&version,
		&createdAt,
		&updatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Person{}, foundationrepository.ErrPersonNotFound
	}
	if err != nil {
		return entity.Person{}, err
	}

	return entity.NewPerson(
		entity.GID(gidValue),
		entity.PersonStatus(status),
		version,
		createdAt,
		updatedAt,
	), nil
}

func (repository *PersonRepository) Update(ctx context.Context, person entity.Person) error {
	result, err := repository.pool.Exec(
		ctx,
		`UPDATE foundation.persons
		 SET status = $2,
		     version = $3,
		     updated_at = $4
		 WHERE gid = $1`,
		string(person.GetGID()),
		string(person.GetStatus()),
		person.GetVersion(),
		person.GetUpdatedAt(),
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return foundationrepository.ErrPersonNotFound
	}

	return nil
}

func (repository *PersonRepository) Delete(ctx context.Context, gid entity.GID) error {
	result, err := repository.pool.Exec(
		ctx,
		`DELETE FROM foundation.persons
		 WHERE gid = $1`,
		string(gid),
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return foundationrepository.ErrPersonNotFound
	}

	return nil
}
