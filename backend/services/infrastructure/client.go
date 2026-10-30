package infrastructure

import (
	"context"

	infrastructureservices "grovehub/services/infrastructure/services"
)

type Address = infrastructureservices.Address
type Message = infrastructureservices.Message

type MailClient interface {
	Send(ctx context.Context, message Message) error
	Health(ctx context.Context) error
}

type mailClient struct {
	smtpService infrastructureservices.SMTPService
}

func NewMailClient(smtpService infrastructureservices.SMTPService) MailClient {
	return &mailClient{smtpService: smtpService}
}

func (client *mailClient) Send(ctx context.Context, message Message) error {
	return client.smtpService.Send(ctx, message)
}

func (client *mailClient) Health(ctx context.Context) error {
	return client.smtpService.Health(ctx)
}
