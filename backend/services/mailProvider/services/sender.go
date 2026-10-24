package services

import (
	"context"

	"grovehub/services/mailProvider/entity"
)

type RegisterSenderInput struct {
	DomainId    string
	Email       string
	DisplayName string
}

type SenderService interface {
	Register(ctx context.Context, input RegisterSenderInput) (entity.Sender, error)
	GetById(ctx context.Context, senderId string) (entity.Sender, error)
	Deactivate(ctx context.Context, senderId string) (entity.Sender, error)
}
