package helpers

import (
	"strings"

	coreerrors "grovehub/libs/core/errors"
	corenormalize "grovehub/libs/core/normalize"
	identityerror "grovehub/services/identity/constants/error"
	"grovehub/services/identity/entity"
)

func NormalizeAccountId(value string) (entity.AccountId, error) {
	normalized, err := corenormalize.UUIDv7(value)
	if err != nil {
		return "", invalidAccount()
	}
	return entity.AccountId(normalized), nil
}

func NormalizePersonGid(value string) (entity.PersonGid, error) {
	normalized, err := corenormalize.UUIDv7(value)
	if err != nil {
		return "", invalidAccount()
	}
	return entity.PersonGid(normalized), nil
}

func NormalizeEmail(value string) (string, error) {
	normalized, err := corenormalize.Email(value)
	if err != nil {
		return "", invalidAccount()
	}
	return normalized, nil
}

func NormalizePhone(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", nil
	}

	normalized, err := corenormalize.Phone(value)
	if err != nil {
		return "", invalidAccount()
	}
	return normalized, nil
}

func invalidAccount() error {
	return coreerrors.NewError(identityerror.InvalidAccount)
}
