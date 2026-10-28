package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/mailProvider/entity"
	mailproviderrepository "grovehub/services/mailProvider/repository"
)

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
	return repository.db.WithContext(ctx).Create(&message).Error
}

func (repository *MessageRepository) FindById(ctx context.Context, messageId entity.MessageId) (entity.Message, error) {
	var message entity.Message
	if err := repository.db.WithContext(ctx).
		Where(&entity.Message{MessageId: messageId}).
		Take(&message).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Message{}, mailproviderrepository.ErrMessageNotFound
		}
		return entity.Message{}, err
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

func (repository *MessageRepository) Update(ctx context.Context, message entity.Message) error {
	result := repository.db.WithContext(ctx).
		Model(&entity.Message{MessageId: message.GetMessageId()}).
		Select("Status", "FailureMessage", "Version", "UpdatedAt").
		Updates(&message)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return mailproviderrepository.ErrMessageNotFound
	}
	return nil
}
