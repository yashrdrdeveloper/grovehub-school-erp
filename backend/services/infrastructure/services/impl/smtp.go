package impl

import (
	"context"
	"net"
	"strconv"
	"strings"

	mail "github.com/wneessen/go-mail"

	infrastructureconfig "grovehub/services/infrastructure/config"
	infrastructurerequest "grovehub/services/infrastructure/dto/request"
	infrastructureservices "grovehub/services/infrastructure/services"
	servicehelpers "grovehub/services/infrastructure/services/helpers"
)

type SMTPService struct {
	config infrastructureconfig.SMTP
}

func NewSMTPService(config infrastructureconfig.SMTP) infrastructureservices.SMTPService {
	return &SMTPService{config: config}
}

func (service *SMTPService) Send(ctx context.Context, message infrastructurerequest.MailMessage) error {
	client, err := mail.NewClient(
		service.config.Host,
		mail.WithPort(service.config.Port),
		mail.WithTLSPolicy(mail.NoTLS),
	)
	if err != nil {
		return err
	}

	mailMessage := mail.NewMsg()
	if err := mailMessage.FromFormat(message.From.DisplayName, message.From.Email); err != nil {
		return err
	}
	for _, recipient := range message.To {
		if err := mailMessage.AddTo(recipient); err != nil {
			return err
		}
	}
	for _, recipient := range message.Cc {
		if err := mailMessage.AddCc(recipient); err != nil {
			return err
		}
	}
	for _, recipient := range message.Bcc {
		if err := mailMessage.AddBcc(recipient); err != nil {
			return err
		}
	}

	internetMessageID := strings.TrimSpace(message.InternetMessageID)
	if internetMessageID != "" {
		mailMessage.SetMessageIDWithValue(servicehelpers.TrimMessageIDBrackets(internetMessageID))
	}
	if inReplyTo := strings.TrimSpace(message.InReplyTo); inReplyTo != "" {
		mailMessage.SetGenHeader(mail.HeaderInReplyTo, servicehelpers.FormatMessageIDHeader(inReplyTo))
	}
	if references := servicehelpers.FormatReferenceMessageIDs(message.ReferenceMessageIDs); references != "" {
		mailMessage.SetGenHeader(mail.HeaderReferences, references)
	}

	mailMessage.Subject(message.Subject)
	if message.TextBody != "" {
		mailMessage.SetBodyString(mail.TypeTextPlain, message.TextBody)
		if message.HTMLBody != "" {
			mailMessage.AddAlternativeString(mail.TypeTextHTML, message.HTMLBody)
		}
	} else {
		mailMessage.SetBodyString(mail.TypeTextHTML, message.HTMLBody)
	}

	sendContext, cancel := context.WithTimeout(ctx, service.config.SubmitTimeout)
	defer cancel()
	return client.DialAndSendWithContext(sendContext, mailMessage)
}

func (service *SMTPService) Health(ctx context.Context) error {
	healthContext, cancel := context.WithTimeout(ctx, service.config.HealthTimeout)
	defer cancel()

	connection, err := (&net.Dialer{}).DialContext(
		healthContext,
		"tcp",
		net.JoinHostPort(service.config.Host, strconv.Itoa(service.config.Port)),
	)
	if err != nil {
		return err
	}
	defer func() { _ = connection.Close() }()
	return nil
}
