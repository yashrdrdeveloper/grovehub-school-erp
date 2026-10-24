package helpers

import (
	"strings"
	"time"

	coresearch "grovehub/libs/core/search"
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

func contactKey(contactType entity.ContactType, normalizedValue string) string {
	return string(contactType) + "\x00" + normalizedValue
}
