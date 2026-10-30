package mail

import "context"

type Address struct {
	Email       string
	DisplayName string
}

type Message struct {
	From                Address
	To                  []string
	Cc                  []string
	Bcc                 []string
	Subject             string
	TextBody            string
	HTMLBody            string
	InternetMessageID   string
	InReplyTo           string
	ReferenceMessageIDs []string
}

type MailClient interface {
	Send(ctx context.Context, message Message) error
	Health(ctx context.Context) error
}
