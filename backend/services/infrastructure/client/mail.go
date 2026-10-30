package client

import (
	"context"

	infrastructurerequest "grovehub/services/infrastructure/dto/request"
)

type Address = infrastructurerequest.MailAddress
type Message = infrastructurerequest.MailMessage

type MailClient interface {
	Send(ctx context.Context, message Message) error
	Health(ctx context.Context) error
}
