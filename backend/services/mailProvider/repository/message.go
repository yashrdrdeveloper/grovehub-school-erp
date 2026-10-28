package repository

import (
	"context"
	"errors"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/mailProvider/entity"
)

var ErrMessageNotFound = errors.New("mail provider message not found")

type MessageRepository interface {
	Create(ctx context.Context, message entity.Message) error
	FindById(ctx context.Context, messageId entity.MessageId) (entity.Message, error)
	Search(ctx context.Context, query coresearch.Query[entity.MessageId]) (coresearch.QueryResult[entity.Message], error)
	Update(ctx context.Context, message entity.Message) error
}
