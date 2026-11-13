package helpers

import (
    "sort"
    "strings"

    "grovehub/services/foundation/entity"
)

// BuildPersonSearchDetails creates the aggregate search projection from Person
// and Contact values; the Person repository persists it in foundation.persons.
func BuildPersonSearchDetails(person entity.Person, personContacts []entity.Contact) string {
	contacts := CopyContacts(personContacts)
	sort.Slice(contacts, func(i, j int) bool {
		if contacts[i].Type == contacts[j].Type {
			return contacts[i].NormalizedValue < contacts[j].NormalizedValue
		}
		return contacts[i].Type < contacts[j].Type
	})
	parts := make([]string, 0, 4+len(contacts)*2)
	parts = append(parts, person.GID, person.FirstName)
	if person.MiddleName != nil { parts = append(parts, *person.MiddleName) }
	parts = append(parts, person.LastName)
	for _, contact := range contacts { parts = append(parts, contact.Value, contact.NormalizedValue) }
	return strings.ToLower(strings.TrimSpace(strings.Join(parts, " ")))
}
