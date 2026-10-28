package client

import (
	"context"

	mailproviderresponse "grovehub/services/mailProvider/dto/response"
)

type SendRequest struct {
	IdempotencyKey string
	GID            string
	SenderId       string
	To             []string
	Cc             []string
	Bcc            []string
	Subject        string
	TextBody       string
	HTMLBody       string
}

type Client interface {
	Send(ctx context.Context, request SendRequest) (mailproviderresponse.Message, error)
	Retry(ctx context.Context, messageId string) (mailproviderresponse.Message, error)
}
