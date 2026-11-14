// Package helpers contains helpers owned exclusively by the Person Create use case.
package helpers

import (
    "time"

    corehelpers "grovehub/libs/core/helpers"
    persondto "grovehub/services/foundation/client/person/dto"
    "grovehub/services/foundation/internal/app/entity"
)

// NewPerson is use-case construction, not entity behavior.
func NewPerson(gid string, input persondto.CreatePerson, now time.Time) entity.Person {
	return entity.Person{
		GID: gid,
		FirstName: input.FirstName,
		MiddleName: corehelpers.OptionalString(input.MiddleName),
		LastName: input.LastName,
		Status: entity.PersonStatusActive,
		Version: 1,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
