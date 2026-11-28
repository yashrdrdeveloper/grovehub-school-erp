package helpers

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils"

func ValidID(value string) bool {
	return utils.ValidateUUID(value) == nil
}
