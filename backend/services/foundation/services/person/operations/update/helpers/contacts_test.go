package helpers

import (
    "testing"

    "grovehub/services/foundation/entity"
)

func TestContactsEqualUsesIdentityAndUserFacingFields(t *testing.T) {
    original := []entity.Contact{{Type: entity.ContactTypeEmail, NormalizedValue: "p@example.com", Value: "p@example.com", IsPrimary: true, IsVerified: true}}
    duplicate := []entity.Contact{{Type: entity.ContactTypeEmail, NormalizedValue: "p@example.com", Value: "p@example.com", IsPrimary: true, IsVerified: true}}
    if !ContactsEqual(original, duplicate) { t.Fatal("equivalent contacts should compare equal") }
    if !ContactsEqual(nil, nil) { t.Fatal("empty contact slices must compare equal") }
    duplicate[0].IsVerified = false
    if ContactsEqual(original, duplicate) { t.Fatal("verification changes must not be ignored") }
    duplicate[0].IsVerified = true
    duplicate[0].Type = entity.ContactTypePhone
    if ContactsEqual(original, duplicate) { t.Fatal("different contact types must not match") }
}
