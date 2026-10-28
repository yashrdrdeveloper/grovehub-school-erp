package helpers

import (
	"strings"

	coreerrors "grovehub/libs/core/errors"
	corenormalize "grovehub/libs/core/normalize"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
	"grovehub/services/mailProvider/entity"
)

func NormalizeDomain(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(value), "."))
	if normalized == "" || len(normalized) > 253 || strings.Contains(normalized, "..") {
		return "", coreerrors.NewError(mailprovidererror.InvalidDomain)
	}

	labels := strings.Split(normalized, ".")
	if len(labels) < 2 {
		return "", coreerrors.NewError(mailprovidererror.InvalidDomain)
	}
	for _, label := range labels {
		if !validDomainLabel(label) {
			return "", coreerrors.NewError(mailprovidererror.InvalidDomain)
		}
	}
	return normalized, nil
}

func NormalizeDomainId(value string) (entity.DomainId, error) {
	normalized, err := corenormalize.UUIDv7(value)
	if err != nil {
		return "", coreerrors.NewError(mailprovidererror.InvalidDomain)
	}
	return entity.DomainId(normalized), nil
}

func validDomainLabel(label string) bool {
	if label == "" || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
		return false
	}
	for _, character := range label {
		if (character >= 'a' && character <= 'z') ||
			(character >= '0' && character <= '9') || character == '-' {
			continue
		}
		return false
	}
	return true
}
