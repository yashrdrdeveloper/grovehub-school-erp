package helpers

import (
    "strings"
    "time"

    "grovehub/services/foundation/internal/app/entity"
)

type cursorValueExtractor func(entity.Person) string

var personCursorExtractors = map[string]cursorValueExtractor{
	"gid":       func(person entity.Person) string { return person.GID },
	"firstName": func(person entity.Person) string { return strings.ToLower(person.FirstName) },
	"lastName":  func(person entity.Person) string { return strings.ToLower(person.LastName) },
	"createdAt": func(person entity.Person) string { return person.CreatedAt.UTC().Format(time.RFC3339Nano) },
	"updatedAt": func(person entity.Person) string { return person.UpdatedAt.UTC().Format(time.RFC3339Nano) },
}

func personCursorValue(person entity.Person, field string) string {
	extract, ok := personCursorExtractors[field]
	if !ok {
		return ""
	}
	return extract(person)
}
