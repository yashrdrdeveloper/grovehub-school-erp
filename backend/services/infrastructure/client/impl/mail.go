package impl

import (
	"context"

	infrastructureclient "grovehub/services/infrastructure/client"
	infrastructureservices "grovehub/services/infrastructure/services"
)

type MailClient struct {
	smtpService infrastructureservices.SMTPService
}

func NewMailClient(smtpService infrastructureservices.SMTPService) infrastructureclient.MailClient {
	return &MailClient{smtpService: smtpService}
}

func (client *MailClient) Send(ctx context.Context, message infrastructureclient.Message) error {
	return client.smtpService.Send(ctx, message)
}

func (client *MailClient) Health(ctx context.Context) error {
	return client.smtpService.Health(ctx)
}
