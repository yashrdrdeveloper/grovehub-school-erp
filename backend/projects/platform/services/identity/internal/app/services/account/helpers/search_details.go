package helpers

import (
	"strings"

	"grovehub/services/identity/internal/app/entity"
)

// SearchDetails contains only the indexed search projection, not user secrets.
func SearchDetails(account entity.Account) string {
	username := ""
	if account.Username != nil {
		username = *account.Username
	}
	return strings.ToLower(strings.Join([]string{account.Email, username}, " "))
}
