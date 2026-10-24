package impl

import (
	"context"

	mailproviderclient "grovehub/services/mailProvider/client"
	mailproviderresponse "grovehub/services/mailProvider/dto/response"
	mailprovidermappers "grovehub/services/mailProvider/mappers"
	mailproviderservices "grovehub/services/mailProvider/services"
)

type Client struct {
	service mailproviderservices.MessageService
	mapper  mailprovidermappers.MessageMapper
}

func NewClient(service mailproviderservices.MessageService, mapper mailprovidermappers.MessageMapper) mailproviderclient.Client {
	return &Client{service: service, mapper: mapper}
}

func (client *Client) Send(ctx context.Context, request mailproviderclient.SendRequest) (mailproviderresponse.Message, error) {
	message, err := client.service.Send(ctx, mailproviderservices.SendMessageInput{
		SenderId: request.SenderId, To: request.To, Cc: request.Cc, Bcc: request.Bcc,
		Subject: request.Subject, TextBody: request.TextBody, HTMLBody: request.HTMLBody,
	})
	if err != nil {
		return mailproviderresponse.Message{}, err
	}
	return client.mapper.ToResponse(message), nil
}
