package services

import (
	"context"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/mailProvider/entity"
)

type SendMessageInput struct {
	IdempotencyKey string
	GID            string
	SenderId       string
	To             []string
	Cc             []string
	Bcc            []string
	Subject        string
	TextBody       string
	HTMLBody       string
}

type ReplyMessageInput struct {
	IdempotencyKey string
	TextBody       string
	HTMLBody       string
}

type MessageService interface {
	Send(ctx context.Context, input SendMessageInput) (entity.Message, error)
	Retry(ctx context.Context, messageId string) (entity.Message, error)
	Resend(ctx context.Context, messageId string, idempotencyKey string) (entity.Message, error)
	Reply(ctx context.Context, messageId string, input ReplyMessageInput) (entity.Message, error)
	GetById(ctx context.Context, messageId string) (entity.Message, error)
	Search(ctx context.Context, request coresearch.Request) (coresearch.Result[entity.Message], error)
}
