package services

import (
	"context"

	"grovehub/services/mailProvider/entity"
)

type SendMessageInput struct {
	SenderId string
	To       []string
	Cc       []string
	Bcc      []string
	Subject  string
	TextBody string
	HTMLBody string
}

type MessageService interface {
	Send(ctx context.Context, input SendMessageInput) (entity.Message, error)
	GetById(ctx context.Context, messageId string) (entity.Message, error)
}
