package impl

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	mail "github.com/wneessen/go-mail"

	inframail "grovehub/services/infrastructure/mail"
)

const (
	defaultSMTPHost      = "127.0.0.1"
	defaultSMTPPort      = 2525
	defaultSubmitTimeout = 15 * time.Second
	defaultHealthTimeout = 3 * time.Second
)

type smtpClient struct {
	host string
	port int
}

func NewMailClientFromEnvironment() inframail.MailClient {
	host := firstNonEmpty(
		os.Getenv("MAIL_INFRASTRUCTURE_SMTP_HOST"),
		os.Getenv("MAIL_PROVIDER_SMTP_HOST"),
		defaultSMTPHost,
	)
	port := envPort("MAIL_INFRASTRUCTURE_SMTP_PORT", "MAIL_PROVIDER_SMTP_PORT", defaultSMTPPort)
	return &smtpClient{host: host, port: port}
}

func (client *smtpClient) Send(ctx context.Context, message inframail.Message) error {
	smtpClient, err := mail.NewClient(client.host, mail.WithPort(client.port), mail.WithTLSPolicy(mail.NoTLS))
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

	internetMessageId := strings.TrimSpace(message.InternetMessageID)
	if internetMessageId != "" {
		mailMessage.SetMessageIDWithValue(trimMessageIDBrackets(internetMessageId))
	}
	if inReplyTo := strings.TrimSpace(message.InReplyTo); inReplyTo != "" {
		mailMessage.SetGenHeader(mail.HeaderInReplyTo, formatMessageIDHeader(inReplyTo))
	}
	if len(message.ReferenceMessageIDs) > 0 {
		formattedReferences := make([]string, 0, len(message.ReferenceMessageIDs))
		for _, reference := range message.ReferenceMessageIDs {
			if strings.TrimSpace(reference) == "" {
				continue
			}
			formattedReferences = append(formattedReferences, formatMessageIDHeader(reference))
		}
		if len(formattedReferences) > 0 {
			mailMessage.SetGenHeader(mail.HeaderReferences, strings.Join(formattedReferences, " "))
		}
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

	sendContext, cancel := context.WithTimeout(ctx, defaultSubmitTimeout)
	defer cancel()
	return smtpClient.DialAndSendWithContext(sendContext, mailMessage)
}

func (client *smtpClient) Health(ctx context.Context) error {
	healthContext, cancel := context.WithTimeout(ctx, defaultHealthTimeout)
	defer cancel()

	connection, err := (&net.Dialer{}).DialContext(
		healthContext,
		"tcp",
		net.JoinHostPort(client.host, strconv.Itoa(client.port)),
	)
	if err != nil {
		return err
	}
	defer func() { _ = connection.Close() }()
	return nil
}

func envPort(primary string, fallback string, defaultValue int) int {
	for _, name := range []string{primary, fallback} {
		value := strings.TrimSpace(os.Getenv(name))
		if value == "" {
			continue
		}
		parsed, err := strconv.Atoi(value)
		if err == nil && parsed > 0 && parsed <= 65535 {
			return parsed
		}
	}
	return defaultValue
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func trimMessageIDBrackets(value string) string {
	return strings.Trim(strings.TrimSpace(value), "<>")
}

func formatMessageIDHeader(value string) string {
	return fmt.Sprintf("<%s>", trimMessageIDBrackets(value))
}
