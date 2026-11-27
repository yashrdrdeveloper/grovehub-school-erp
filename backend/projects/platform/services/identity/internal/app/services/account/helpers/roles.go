package helpers

import (
	"encoding/json"

	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
)

func ApplyRoleAssignment(account *entity.Account, assignment accountdto.RoleAssignment) error {
	secondary, err := json.Marshal(assignment.SecondaryRoleIDs)
	if err != nil {
		return err
	}
	if string(secondary) == "null" {
		secondary = []byte("[]")
	}
	primary := assignment.PrimaryRoleID
	account.PrimaryRoleID = &primary
	account.SecondaryRoleIDs = secondary
	return nil
}
