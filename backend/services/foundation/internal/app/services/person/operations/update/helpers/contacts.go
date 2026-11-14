package helpers

import "grovehub/services/foundation/internal/app/entity"

// contactIdentity compares normalized contact identities without audit fields.
type contactIdentity struct {
    Type entity.ContactType
    NormalizedValue string
}

// ContactsEqual compares PATCH contact values without considering audit timestamps.
func ContactsEqual(left, right []entity.Contact) bool {
	if len(left) != len(right) {
		return false
	}
	rightByKey := make(map[contactIdentity]entity.Contact, len(right))
	for _, contact := range right {
		rightByKey[contactIdentity{Type: contact.Type, NormalizedValue: contact.NormalizedValue}] = contact
	}
	for _, contact := range left {
		other, exists := rightByKey[contactIdentity{Type: contact.Type, NormalizedValue: contact.NormalizedValue}]
		if !exists ||
			other.Value != contact.Value ||
			other.IsPrimary != contact.IsPrimary ||
			other.IsVerified != contact.IsVerified {
			return false
		}
	}
	return true
}
