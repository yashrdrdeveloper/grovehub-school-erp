package transport

import (
	"context"

	inframail "grovehub/services/infrastructure/mail"
	"grovehub/services/mailProvider/entity"
)

type Sender interface {
	Send(ctx context.Context, message entity.Message, sender entity.Sender) error
}

type sender struct {
	mailClient inframail.MailClient
}

func NewSender(mailClient inframail.MailClient) Sender {
	return &sender{mailClient: mailClient}
}

func (adapter *sender) Send(ctx context.Context, message entity.Message, from entity.Sender) error {
	return adapter.mailClient.Send(ctx, inframail.Message{
		From: inframail.Address{
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
