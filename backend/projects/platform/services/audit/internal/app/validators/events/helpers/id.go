package helpers

import "grovehub/libs/core/utils"

func ValidID(value string) bool {
	return utils.ValidateUUID(value) == nil
}
