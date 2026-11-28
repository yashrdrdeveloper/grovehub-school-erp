package impl

import (
	client "grovehub/services/integration/client/mail"
	"grovehub/services/integration/internal/app/services/mail"
)

type mailClient struct {
	service mail.Service
}

func NewMailClient(service mail.Service) client.Client {
	return &mailClient{service: service}
}
