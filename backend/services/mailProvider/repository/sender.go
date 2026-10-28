package repository

import (
	"context"
	"errors"

	"grovehub/services/mailProvider/entity"
)

var (
	ErrSenderNotFound      = errors.New("mail provider sender not found")
	ErrSenderAlreadyExists = errors.New("mail provider sender already exists")
)

type SenderRepository interface {
	Create(ctx context.Context, sender entity.Sender) error
	FindAll(ctx context.Context) ([]entity.Sender, error)
	FindById(ctx context.Context, senderId entity.SenderId) (entity.Sender, error)
	FindByEmail(ctx context.Context, email string) (entity.Sender, error)
	Update(ctx context.Context, sender entity.Sender) error
}
