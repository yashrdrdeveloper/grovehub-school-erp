package validator_test

import (
	"testing"

	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/services/account/operations/profile/helpers"
)

func TestProfilePatchDoesNotReplaceContacts(t *testing.T) {
	patch, err := helpers.BuildPatch(accountdto.UpdateProfile{FirstName: " Rahul ", LastName: " Sharma "})
	if err != nil {
		t.Fatal(err)
	}
	if patch.Contacts.Set || !patch.FirstName.Set || !patch.LastName.Set || !patch.MiddleName.Set || patch.MiddleName.Value != nil {
		t.Fatal("profile must update names only and preserve contacts")
	}
	if *patch.FirstName.Value != "Rahul" || *patch.LastName.Value != "Sharma" {
		t.Fatal("names not normalized")
	}
	if _, err := helpers.BuildPatch(accountdto.UpdateProfile{FirstName: "", LastName: "Sharma"}); err == nil {
		t.Fatal("empty name accepted")
	}
}
