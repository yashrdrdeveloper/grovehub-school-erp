package helpers

import (
	coresearch "grovehub/libs/core/search"
	"grovehub/libs/core/utils"
	"grovehub/services/audit/internal/app/services/activity"
	"grovehub/services/audit/internal/app/validators/activity/constants"
)

func ValidateID(value string, invalid error) error {
	if utils.ValidateUUID(value) != nil {
		return invalid
	}
	return nil
}

func NormalizeID(_ coresearch.Operator, value string) (string, error) {
	return value, ValidateID(value, activity.ErrInvalidSearch)
}

func NormalizeText(_ coresearch.Operator, value string) (string, error) {
	normalized, err := utils.NormalizeText(value, constants.MaxFilterTextBytes)
	if err != nil {
		return "", activity.ErrInvalidSearch
	}
	return normalized, nil
}

func NormalizeTime(_ coresearch.Operator, value string) (string, error) {
	normalized, err := utils.NormalizeStrictRFC3339UTC(value)
	if err != nil {
		return "", activity.ErrInvalidSearch
	}
	return normalized, nil
}
