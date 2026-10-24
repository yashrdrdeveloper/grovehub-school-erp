package impl

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"grovehub/services/mailProvider/entity"
	mailproviderrepository "grovehub/services/mailProvider/repository"
)

type MessageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) mailproviderrepository.MessageRepository {
	return &MessageRepository{pool: pool}
}

func (repository *MessageRepository) Create(ctx context.Context, message entity.Message) error {
	_, err := repository.pool.Exec(ctx, `INSERT INTO mail_provider.messages (
		message_id, sender_id, to_addresses, cc_addresses, bcc_addresses, subject,
		text_body, html_body, status, failure_message, version, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NULLIF($10, ''), $11, $12, $13)`,
		string(message.GetMessageId()), string(message.GetSenderId()), message.GetTo(), message.GetCc(), message.GetBcc(),
		message.GetSubject(), message.GetTextBody(), message.GetHTMLBody(), string(message.GetStatus()), message.GetFailureMessage(),
		message.GetVersion(), message.GetCreatedAt(), message.GetUpdatedAt(),
	)
	return err
}

func (repository *MessageRepository) FindById(ctx context.Context, messageId entity.MessageId) (entity.Message, error) {
	return scanMessage(repository.pool.QueryRow(ctx, messageSelect+` WHERE message_id = $1`, string(messageId)))
}

func (repository *MessageRepository) Update(ctx context.Context, message entity.Message) error {
	result, err := repository.pool.Exec(ctx, `UPDATE mail_provider.messages
	SET status = $2, failure_message = NULLIF($3, ''), version = $4, updated_at = $5 WHERE message_id = $1`,
		string(message.GetMessageId()), string(message.GetStatus()), message.GetFailureMessage(), message.GetVersion(), message.GetUpdatedAt(),
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return mailproviderrepository.ErrMessageNotFound
	}
	return nil
}

const messageSelect = `SELECT message_id, sender_id, to_addresses, cc_addresses, bcc_addresses, subject,
	text_body, html_body, status, failure_message, version, created_at, updated_at FROM mail_provider.messages`

func scanMessage(row pgx.Row) (entity.Message, error) {
	var messageId, senderId, subject, textBody, htmlBody, status string
	var failureMessage *string
	var to, cc, bcc []string
	var version int64
	var createdAt, updatedAt time.Time
	if err := row.Scan(&messageId, &senderId, &to, &cc, &bcc, &subject, &textBody, &htmlBody, &status, &failureMessage, &version, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Message{}, mailproviderrepository.ErrMessageNotFound
		}
		return entity.Message{}, err
	}
	failure := ""
	if failureMessage != nil {
		failure = *failureMessage
	}
	return entity.NewMessage(entity.MessageId(messageId), entity.SenderId(senderId), to, cc, bcc, subject, textBody, htmlBody,
		entity.MessageStatus(status), failure, version, createdAt, updatedAt), nil
}
