package impl

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	mail "github.com/wneessen/go-mail"

	"grovehub/services/mailProvider/entity"
	mailprovidertransport "grovehub/services/mailProvider/transport"
)

const (
	defaultHarakaHost    = "127.0.0.1"
	defaultHarakaPort    = 2525
	defaultSubmitTimeout = 15 * time.Second
)

type HarakaSender struct {
	host string
	port int
}

func NewHarakaSenderFromEnvironment() mailprovidertransport.Sender {
	host := os.Getenv("MAIL_PROVIDER_SMTP_HOST")
	if host == "" {
		host = defaultHarakaHost
	}
	port := defaultHarakaPort
	if value := os.Getenv("MAIL_PROVIDER_SMTP_PORT"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 && parsed <= 65535 {
			port = parsed
		}
	}
	return &HarakaSender{host: host, port: port}
}

func (sender *HarakaSender) Send(ctx context.Context, message entity.Message, from entity.Sender) error {
	client, err := mail.NewClient(sender.host, mail.WithPort(sender.port), mail.WithTLSPolicy(mail.NoTLS))
	if err != nil {
		return err
	}

	mailMessage := mail.NewMsg()
	if err := mailMessage.FromFormat(from.GetDisplayName(), from.GetEmail()); err != nil {
		return err
	}
	for _, recipient := range message.GetTo() {
		if err := mailMessage.AddTo(recipient); err != nil {
			return err
		}
	}
	for _, recipient := range message.GetCc() {
		if err := mailMessage.AddCc(recipient); err != nil {
			return err
		}
	}
	for _, recipient := range message.GetBcc() {
		if err := mailMessage.AddBcc(recipient); err != nil {
			return err
		}
	}

	internetMessageId := strings.TrimSpace(message.GetInternetMessageId())
	if internetMessageId != "" {
		mailMessage.SetMessageIDWithValue(trimMessageIDBrackets(internetMessageId))
	}
	if inReplyTo := strings.TrimSpace(message.GetInReplyTo()); inReplyTo != "" {
		mailMessage.SetGenHeader(mail.HeaderInReplyTo, formatMessageIDHeader(inReplyTo))
	}
	if references := message.GetReferenceMessageIds(); len(references) > 0 {
		formattedReferences := make([]string, 0, len(references))
		for _, reference := range references {
			if strings.TrimSpace(reference) == "" {
				continue
			}
			formattedReferences = append(formattedReferences, formatMessageIDHeader(reference))
		}
		if len(formattedReferences) > 0 {
			mailMessage.SetGenHeader(mail.HeaderReferences, strings.Join(formattedReferences, " "))
		}
	}

	mailMessage.Subject(message.GetSubject())
	if message.GetTextBody() != "" {
		mailMessage.SetBodyString(mail.TypeTextPlain, message.GetTextBody())
		if message.GetHTMLBody() != "" {
			mailMessage.AddAlternativeString(mail.TypeTextHTML, message.GetHTMLBody())
		}
	} else {
		mailMessage.SetBodyString(mail.TypeTextHTML, message.GetHTMLBody())
	}

	sendContext, cancel := context.WithTimeout(ctx, defaultSubmitTimeout)
	defer cancel()
	return client.DialAndSendWithContext(sendContext, mailMessage)
}

func trimMessageIDBrackets(value string) string {
	return strings.Trim(strings.TrimSpace(value), "<>")
}

func formatMessageIDHeader(value string) string {
	return fmt.Sprintf("<%s>", trimMessageIDBrackets(value))
}
