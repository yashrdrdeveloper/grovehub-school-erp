package repository

import (
	"context"
	"errors"

	"grovehub/services/mailProvider/entity"
)

var ErrMessageNotFound = errors.New("mail provider message not found")

type MessageRepository interface {
	Create(ctx context.Context, message entity.Message) error
	FindById(ctx context.Context, messageId entity.MessageId) (entity.Message, error)
	Update(ctx context.Context, message entity.Message) error
}
