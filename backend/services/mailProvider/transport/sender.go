package transport

import (
	"context"

	infrastructureclient "grovehub/services/infrastructure/client"
	"grovehub/services/mailProvider/entity"
)

type Sender interface {
	Send(ctx context.Context, message entity.Message, sender entity.Sender) error
}

type sender struct {
	mailClient infrastructureclient.MailClient
}

func NewSender(mailClient infrastructureclient.MailClient) Sender {
	return &sender{mailClient: mailClient}
}

func (adapter *sender) Send(ctx context.Context, message entity.Message, from entity.Sender) error {
	return adapter.mailClient.Send(ctx, infrastructureclient.Message{
		From: infrastructureclient.Address{
			Email:       from.GetEmail(),
			DisplayName: from.GetDisplayName(),
		},
		To:                  message.GetTo(),
		Cc:                  message.GetCc(),
		Bcc:                 message.GetBcc(),
		Subject:             message.GetSubject(),
		TextBody:            message.GetTextBody(),
		HTMLBody:            message.GetHTMLBody(),
		InternetMessageID:   message.GetInternetMessageId(),
		InReplyTo:           message.GetInReplyTo(),
		ReferenceMessageIDs: message.GetReferenceMessageIds(),
	})
}
