package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"grovehub/services/mailProvider/entity"
	mailproviderrepository "grovehub/services/mailProvider/repository"
)

type MessageRepository struct {
	db *gorm.DB
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
