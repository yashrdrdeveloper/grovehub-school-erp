package impl

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/mailProvider/entity"
	mailproviderrepository "grovehub/services/mailProvider/repository"
)

const uniqueMessageIdempotencyKeyConstraint = "uq_mail_provider_messages_idempotency_key"

type MessageRepository struct {
	db *gorm.DB
}

var messageSearchConfig = coresearch.Config{
	SearchExpression: "CONCAT_WS(' ', m.message_id::text, m.gid::text, m.sender_id::text, m.subject, m.status, array_to_string(m.to_addresses, ' '), array_to_string(m.cc_addresses, ' '), array_to_string(m.bcc_addresses, ' '))",
	CursorExpression: "m.message_id",
	ScalarFields: map[string]coresearch.ScalarField{
		"messageId": {Expression: "m.message_id::text"},
		"gid":       {Expression: "m.gid::text"},
		"senderId":  {Expression: "m.sender_id::text"},
		"subject":   {Expression: "LOWER(m.subject)"},
		"status":    {Expression: "m.status"},
		"createdAt": {Expression: "m.created_at", Cast: "::timestamptz"},
		"updatedAt": {Expression: "m.updated_at", Cast: "::timestamptz"},
	},
	SortFields: map[string]coresearch.SortField{
		"messageId": {Expression: "m.message_id", Cast: "::uuid"},
		"gid":       {Expression: "m.gid", Cast: "::uuid"},
		"senderId":  {Expression: "m.sender_id", Cast: "::uuid"},
		"subject":   {Expression: "LOWER(m.subject)", Cast: "::text"},
		"status":    {Expression: "m.status", Cast: "::text"},
		"createdAt": {Expression: "m.created_at", Cast: "::timestamptz"},
		"updatedAt": {Expression: "m.updated_at", Cast: "::timestamptz"},
	},
}

func NewMessageRepository(db *gorm.DB) mailproviderrepository.MessageRepository {
	return &MessageRepository{db: db}
}

func (repository *MessageRepository) Create(ctx context.Context, message entity.Message) error {
	return mapMessageWriteError(repository.db.WithContext(ctx).Create(&message).Error)
}

func (repository *MessageRepository) FindById(ctx context.Context, messageId entity.MessageId) (entity.Message, error) {
	var message entity.Message
	if err := repository.db.WithContext(ctx).
		Where(&entity.Message{MessageId: messageId}).
		Take(&message).Error; err != nil {
		return mapMessageReadError(err)
	}
	return message, nil
}

func (repository *MessageRepository) FindByIdempotencyKey(ctx context.Context, key string) (entity.Message, error) {
	var message entity.Message
	if err := repository.db.WithContext(ctx).
		Where(&entity.Message{IdempotencyKey: key}).
		Take(&message).Error; err != nil {
		return mapMessageReadError(err)
	}
	return message, nil
}

func (repository *MessageRepository) Search(
	ctx context.Context,
	query coresearch.Query[entity.MessageId],
) (coresearch.QueryResult[entity.Message], error) {
	searchDB, err := coresearch.Apply(
		repository.db.WithContext(ctx).
			Model(&entity.Message{}).
			Table(entity.Message{}.TableName()+" AS m").
			Select("m.*"),
		query,
		messageSearchConfig,
	)
	if err != nil {
		return coresearch.QueryResult[entity.Message]{}, err
	}

	messages := make([]entity.Message, 0, query.Limit+1)
	if err := searchDB.Find(&messages).Error; err != nil {
		return coresearch.QueryResult[entity.Message]{}, err
	}

	hasNext := len(messages) > query.Limit
	if hasNext {
		messages = messages[:query.Limit]
	}
	return coresearch.QueryResult[entity.Message]{Items: messages, HasNext: hasNext}, nil
}

func (repository *MessageRepository) PrepareRetry(
	ctx context.Context,
	message entity.Message,
	retryAt time.Time,
) (entity.Message, error) {
	result := repository.db.WithContext(ctx).
		Model(&entity.Message{}).
		Where(&entity.Message{MessageId: message.GetMessageId(), Status: entity.MessageStatusFailed}).
		Where("version = ?", message.GetVersion()).
		Updates(map[string]any{
			"status":                entity.MessageStatusQueued,
			"failure_message":       nil,
			"submission_started_at": nil,
			"version":               message.GetVersion() + 1,
			"updated_at":            retryAt,
		})
	if result.Error != nil {
		return entity.Message{}, result.Error
	}
	if result.RowsAffected == 0 {
		return entity.Message{}, mailproviderrepository.ErrMessageVersionConflict
	}

	message.Status = entity.MessageStatusQueued
	message.FailureMessage = nil
	message.SubmissionStartedAt = nil
	message.Version = message.GetVersion() + 1
	message.UpdatedAt = retryAt
	return message, nil
}

func (repository *MessageRepository) ClaimSubmission(
	ctx context.Context,
	message entity.Message,
	startedAt time.Time,
) (entity.Message, error) {
	result := repository.db.WithContext(ctx).
		Model(&entity.Message{}).
		Where(&entity.Message{MessageId: message.GetMessageId(), Status: entity.MessageStatusQueued}).
		Where("version = ?", message.GetVersion()).
		Where("submission_started_at IS NULL").
		Updates(map[string]any{
			"submission_started_at": startedAt,
			"version":               message.GetVersion() + 1,
			"updated_at":            startedAt,
		})
	if result.Error != nil {
		return entity.Message{}, result.Error
	}
	if result.RowsAffected == 0 {
		return entity.Message{}, mailproviderrepository.ErrMessageSubmissionAlreadyClaimed
	}
	return message.
		SetSubmissionStartedAt(startedAt).
		SetVersion(message.GetVersion() + 1).
		SetUpdatedAt(startedAt), nil
}

func (repository *MessageRepository) Update(ctx context.Context, message entity.Message) error {
	expectedVersion := message.GetVersion() - 1
	result := repository.db.WithContext(ctx).
		Model(&entity.Message{}).
		Where(&entity.Message{MessageId: message.GetMessageId()}).
		Where("version = ?", expectedVersion).
		Select("Status", "FailureMessage", "Version", "UpdatedAt").
		Updates(&message)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return mailproviderrepository.ErrMessageVersionConflict
	}
	return nil
}

func mapMessageReadError(err error) (entity.Message, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Message{}, mailproviderrepository.ErrMessageNotFound
	}
	return entity.Message{}, err
}

func mapMessageWriteError(err error) error {
	if err == nil {
		return nil
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) &&
		postgresError.Code == "23505" &&
		postgresError.ConstraintName == uniqueMessageIdempotencyKeyConstraint {
		return mailproviderrepository.ErrMessageIdempotencyKeyAlreadyExists
	}
	return err
}
