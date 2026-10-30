package transport

import (
	"context"

	"grovehub/services/infrastructure"
	"grovehub/services/mailProvider/entity"
)

type Sender interface {
	Send(ctx context.Context, message entity.Message, sender entity.Sender) error
}

type sender struct {
	mailClient infrastructure.MailClient
}

func NewSender(mailClient infrastructure.MailClient) Sender {
	return &sender{mailClient: mailClient}
}

func (adapter *sender) Send(ctx context.Context, message entity.Message, from entity.Sender) error {
	return adapter.mailClient.Send(ctx, infrastructure.Message{
		From: infrastructure.Address{
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
