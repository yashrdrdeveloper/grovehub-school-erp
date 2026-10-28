package repository

import (
	"context"
	"errors"
	"time"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/mailProvider/entity"
)

var (
	ErrMessageNotFound                    = errors.New("mail provider message not found")
	ErrMessageIdempotencyKeyAlreadyExists = errors.New("mail provider message idempotency key already exists")
	ErrMessageVersionConflict             = errors.New("mail provider message version conflict")
	ErrMessageSubmissionAlreadyClaimed    = errors.New("mail provider message submission already claimed")
)

type MessageRepository interface {
	Create(ctx context.Context, message entity.Message) error
	FindById(ctx context.Context, messageId entity.MessageId) (entity.Message, error)
	FindByIdempotencyKey(ctx context.Context, key string) (entity.Message, error)
	Search(ctx context.Context, query coresearch.Query[entity.MessageId]) (coresearch.QueryResult[entity.Message], error)
	PrepareRetry(ctx context.Context, message entity.Message, retryAt time.Time) (entity.Message, error)
	ClaimSubmission(ctx context.Context, message entity.Message, startedAt time.Time) (entity.Message, error)
	Update(ctx context.Context, message entity.Message) error
}
