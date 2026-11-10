package helpers

import (
	"strings"
	"time"

	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
	"grovehub/services/foundation/entity"
)

type cursorValueExtractor func(entity.Person) string

var personCursorExtractors = map[string]cursorValueExtractor{
	"gid":       func(person entity.Person) string { return person.GID },
	"firstName": func(person entity.Person) string { return strings.ToLower(person.FirstName) },
	"lastName":  func(person entity.Person) string { return strings.ToLower(person.LastName) },
	"createdAt": func(person entity.Person) string { return person.CreatedAt.UTC().Format(time.RFC3339Nano) },
	"updatedAt": func(person entity.Person) string { return person.UpdatedAt.UTC().Format(time.RFC3339Nano) },
}

func BuildContacts(request foundationrequest.Contacts, existing []entity.Contact) []entity.Contact {
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

func ContactsEqual(left, right []entity.Contact) bool {
	if len(left) != len(right) {
		return false
	}
	rightByKey := make(map[string]entity.Contact, len(right))
	for _, contact := range right {
		rightByKey[contactKey(contact.Type, contact.NormalizedValue)] = contact
	}
	for _, contact := range left {
		other, exists := rightByKey[contactKey(contact.Type, contact.NormalizedValue)]
		if !exists ||
			other.Value != contact.Value ||
			other.IsPrimary != contact.IsPrimary ||
			other.IsVerified != contact.IsVerified {
			return false
		}
	}
	return true
}

func PersonCursorValue(person entity.Person, field string) string {
	extract, ok := personCursorExtractors[field]
	if !ok {
		return ""
	}
	return extract(person)
}

func EncodePersonSearchCursor(cursor coresearch.CursorPayload) (string, error) {
	return coresearch.EncodeCursor(cursor)
}

func buildContactGroup(
	contactType entity.ContactType,
	values []foundationrequest.ContactValue,
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

// OptionalString maps an empty optional request string to NULL at the service layer.
func OptionalString(value string) *string {
	if value == "" { return nil }
	return &value
}

// NullableString presents an optional persisted value as the empty string.
func NullableString(value *string) string {
	if value == nil { return "" }
	return *value
}

// CopyContacts preserves the previous service API's defensive-copy semantics.
func CopyContacts(contacts []entity.Contact) []entity.Contact {
	if len(contacts) == 0 { return nil }
	return append([]entity.Contact(nil), contacts...)
}

// NewPerson is use-case construction, not entity behavior.
func NewPerson(gid string, input foundationrequest.CreatePerson, now time.Time) entity.Person {
	return entity.Person{
		GID: gid,
		FirstName: input.FirstName,
		MiddleName: OptionalString(input.MiddleName),
		LastName: input.LastName,
		Status: entity.PersonStatusActive,
		Version: 1,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
