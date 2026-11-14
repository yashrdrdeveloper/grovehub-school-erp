// Package helpers contains Person functions shared across operations.
package helpers

import (
    persondto "grovehub/services/foundation/client/person/dto"
    "grovehub/services/foundation/internal/app/entity"
)

func BuildContacts(request persondto.Contacts, existing []entity.Contact) []entity.Contact {
	existingVerification := make(map[string]bool, len(existing))
	for _, contact := range existing {
		existingVerification[contactKey(contact.Type, contact.NormalizedValue)] = contact.IsVerified
	}

	contacts := make([]entity.Contact, 0, len(request.Emails)+len(request.Phones))
	contacts = append(contacts, buildContactGroup(
		entity.ContactTypeEmail,
		request.Emails,
		existingVerification,
	)...)
	contacts = append(contacts, buildContactGroup(
		entity.ContactTypePhone,
		request.Phones,
		existingVerification,
	)...)
	return contacts
}

func buildContactGroup(
	contactType entity.ContactType,
	values []persondto.ContactValue,
	existingVerification map[string]bool,
) []entity.Contact {
	contacts := make([]entity.Contact, 0, len(values))
	for _, contact := range values {
		key := contactKey(contactType, contact.Value)
		contacts = append(contacts, entity.Contact{
			Type: contactType,
			Value: contact.Value,
			NormalizedValue: contact.Value,
			IsPrimary: contact.IsPrimary,
			IsVerified: existingVerification[key],
		})
	}
	return contacts
}

func contactKey(contactType entity.ContactType, normalizedValue string) string {
	return string(contactType) + "\x00" + normalizedValue
}

// AssociateContacts applies Person ownership and audit timestamps to separate Contact entities.
func AssociateContacts(person entity.Person, contacts []entity.Contact) []entity.Contact {
	associated := CopyContacts(contacts)
	for index := range associated {
		associated[index].PersonGID = person.GID
		if associated[index].CreatedAt.IsZero() {
			associated[index].CreatedAt = person.UpdatedAt
		}
		associated[index].UpdatedAt = person.UpdatedAt
	}
	return associated
}

// CopyContacts preserves the previous service API's defensive-copy semantics.
func CopyContacts(contacts []entity.Contact) []entity.Contact {
	if len(contacts) == 0 { return nil }
	return append([]entity.Contact(nil), contacts...)
}
