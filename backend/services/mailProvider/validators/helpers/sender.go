package helpers

import (
	"strings"

	coreerrors "grovehub/libs/core/errors"
	corenormalize "grovehub/libs/core/helpers/normalize"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
	"grovehub/services/mailProvider/entity"
)

func NormalizeSenderId(value string) (entity.SenderId, error) {
	normalized, err := corenormalize.UUIDv7(value)
	if err != nil {
		return "", coreerrors.NewError(mailprovidererror.InvalidSender)
	}
	return entity.SenderId(normalized), nil
}

func NormalizeSenderEmail(value string) (string, error) {
	normalized, err := corenormalize.Email(value)
	if err != nil {
		return "", coreerrors.NewError(mailprovidererror.InvalidSender)
	}
	return normalized, nil
}

func NormalizeDisplayName(value string) (string, error) {
	normalized := strings.TrimSpace(value)
	if normalized == "" || len([]rune(normalized)) > 120 {
		return "", coreerrors.NewError(mailprovidererror.InvalidSender)
	}
	return normalized, nil
}
