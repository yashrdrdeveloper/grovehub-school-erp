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

const uniqueSenderEmailConstraint = "uq_mail_provider_senders_email"

type SenderRepository struct {
	pool *pgxpool.Pool
}

func NewSenderRepository(pool *pgxpool.Pool) mailproviderrepository.SenderRepository {
	return &SenderRepository{pool: pool}
}

func (repository *SenderRepository) Create(ctx context.Context, sender entity.Sender) error {
	_, err := repository.pool.Exec(ctx, `INSERT INTO mail_provider.senders (
		sender_id, domain_id, email, display_name, status, version, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(sender.GetSenderId()), string(sender.GetDomainId()), sender.GetEmail(), sender.GetDisplayName(),
		string(sender.GetStatus()), sender.GetVersion(), sender.GetCreatedAt(), sender.GetUpdatedAt(),
	)
	return mapSenderWriteError(err)
}

func (repository *SenderRepository) FindById(ctx context.Context, senderId entity.SenderId) (entity.Sender, error) {
	return scanSender(repository.pool.QueryRow(ctx, senderSelect+` WHERE sender_id = $1`, string(senderId)))
}

func (repository *SenderRepository) FindByEmail(ctx context.Context, email string) (entity.Sender, error) {
	return scanSender(repository.pool.QueryRow(ctx, senderSelect+` WHERE email = $1`, email))
}

func (repository *SenderRepository) Update(ctx context.Context, sender entity.Sender) error {
	result, err := repository.pool.Exec(ctx, `UPDATE mail_provider.senders
	SET status = $2, version = $3, updated_at = $4 WHERE sender_id = $1`,
		string(sender.GetSenderId()), string(sender.GetStatus()), sender.GetVersion(), sender.GetUpdatedAt(),
	)
	if err != nil {
		return mapSenderWriteError(err)
	}
	if result.RowsAffected() == 0 {
		return mailproviderrepository.ErrSenderNotFound
	}
	return nil
}

const senderSelect = `SELECT sender_id, domain_id, email, display_name, status, version, created_at, updated_at FROM mail_provider.senders`

func scanSender(row pgx.Row) (entity.Sender, error) {
	var senderId, domainId, email, displayName, status string
	var version int64
	var createdAt, updatedAt time.Time
	if err := row.Scan(&senderId, &domainId, &email, &displayName, &status, &version, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Sender{}, mailproviderrepository.ErrSenderNotFound
		}
		return entity.Sender{}, err
	}
	return entity.NewSender(entity.SenderId(senderId), entity.DomainId(domainId), email, displayName, entity.SenderStatus(status), version, createdAt, updatedAt), nil
}

func mapSenderWriteError(err error) error {
	if err == nil {
		return nil
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" && postgresError.ConstraintName == uniqueSenderEmailConstraint {
		return mailproviderrepository.ErrSenderAlreadyExists
	}
	return err
}
