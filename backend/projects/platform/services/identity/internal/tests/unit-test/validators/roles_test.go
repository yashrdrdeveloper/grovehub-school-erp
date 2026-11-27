package validator_test

import (
	"reflect"
	"testing"

	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/validators/account/helpers"
)

func TestRoleAssignment(t *testing.T) {
	result, err := helpers.NormalizeRoleAssignment(accountdto.RoleAssignment{
		PrimaryRoleID: "TEACHER",
		SecondaryRoleIDs: []string{"LIBRARIAN", "COORDINATOR"},
	})
	if err != nil || !reflect.DeepEqual(result.SecondaryRoleIDs, []string{"COORDINATOR", "LIBRARIAN"}) {
		t.Fatalf("unexpected normalized assignment: %+v, %v", result, err)
	}
	for _, input := range []accountdto.RoleAssignment{
		{PrimaryRoleID: ""},
		{PrimaryRoleID: "UNKNOWN"},
		{PrimaryRoleID: "TEACHER", SecondaryRoleIDs: []string{"TEACHER"}},
		{PrimaryRoleID: "TEACHER", SecondaryRoleIDs: []string{"LIBRARIAN", "LIBRARIAN"}},
		{PrimaryRoleID: "TEACHER", SecondaryRoleIDs: []string{"UNKNOWN"}},
	} {
		if _, err := helpers.NormalizeRoleAssignment(input); err == nil {
			t.Errorf("invalid assignment accepted: %+v", input)
		}
	}
}
