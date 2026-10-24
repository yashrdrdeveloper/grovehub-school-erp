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
	"gid":       func(person entity.Person) string { return string(person.GetGID()) },
	"firstName": func(person entity.Person) string { return strings.ToLower(person.GetFirstName()) },
	"lastName":  func(person entity.Person) string { return strings.ToLower(person.GetLastName()) },
	"createdAt": func(person entity.Person) string { return person.GetCreatedAt().UTC().Format(time.RFC3339Nano) },
	"updatedAt": func(person entity.Person) string { return person.GetUpdatedAt().UTC().Format(time.RFC3339Nano) },
}

func BuildContacts(request foundationrequest.Contacts, existing []entity.Contact) []entity.Contact {
	existingVerification := make(map[string]bool, len(existing))
	for _, contact := range existing {
		existingVerification[contactKey(contact.GetType(), contact.GetNormalizedValue())] = contact.GetIsVerified()
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
		rightByKey[contactKey(contact.GetType(), contact.GetNormalizedValue())] = contact
	}
	for _, contact := range left {
		other, exists := rightByKey[contactKey(contact.GetType(), contact.GetNormalizedValue())]
		if !exists ||
			other.GetValue() != contact.GetValue() ||
			other.GetIsPrimary() != contact.GetIsPrimary() ||
			other.GetIsVerified() != contact.GetIsVerified() {
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
		contacts = append(contacts, entity.NewContact(
			contactType,
			contact.Value,
			contact.Value,
			contact.IsPrimary,
			existingVerification[key],
		))
	}
	return contacts
}

func contactKey(contactType entity.ContactType, normalizedValue string) string {
	return string(contactType) + "\x00" + normalizedValue
}
