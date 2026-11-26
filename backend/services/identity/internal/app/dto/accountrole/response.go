package accountrole

import "grovehub/services/identity/internal/app/entity"

type Assignment struct {
	AccountID string `json:"accountId"`
	Status entity.AccountStatus `json:"status"`
	Version int64 `json:"version"`
	PrimaryRoleID string `json:"primaryRoleId"`
	SecondaryRoleIDs []string `json:"secondaryRoleIds"`
}
