package helpers

import (
	"strings"

	coreerrors "grovehub/libs/core/errors"
	corenormalize "grovehub/libs/core/helpers/normalize"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
	"grovehub/services/mailProvider/entity"
)

const maxRecipients = 100

func NormalizeMessageId(value string) (entity.MessageId, error) {
	normalized, err := corenormalize.UUIDv7(value)
	if err != nil {
		return "", coreerrors.NewError(mailprovidererror.InvalidMessage)
	}
	return entity.MessageId(normalized), nil
}

func NormalizeRecipients(values []string) ([]string, error) {
	normalized := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		email, err := corenormalize.Email(value)
		if err != nil {
			return nil, coreerrors.NewError(mailprovidererror.InvalidMessage)
		}
		if _, ok := seen[email]; ok {
			continue
		}
		seen[email] = struct{}{}
		normalized = append(normalized, email)
	}
	return normalized, nil
}

func ValidateMessageContent(to, cc, bcc []string, subject, textBody, htmlBody string) error {
	if len(to) == 0 || len(to)+len(cc)+len(bcc) > maxRecipients {
		return coreerrors.NewError(mailprovidererror.InvalidMessage)
	}
	if strings.TrimSpace(subject) == "" || len([]rune(subject)) > 998 {
		return coreerrors.NewError(mailprovidererror.InvalidMessage)
	}
	if strings.TrimSpace(textBody) == "" && strings.TrimSpace(htmlBody) == "" {
		return coreerrors.NewError(mailprovidererror.InvalidMessage)
	}
	return nil
}
