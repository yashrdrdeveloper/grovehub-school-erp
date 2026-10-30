package services

import (
	"context"

	infrastructurerequest "grovehub/services/infrastructure/dto/request"
)

type SMTPService interface {
	Send(ctx context.Context, message infrastructurerequest.MailMessage) error
	Health(ctx context.Context) error
}
