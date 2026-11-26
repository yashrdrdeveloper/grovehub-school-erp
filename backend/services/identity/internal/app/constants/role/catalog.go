package role

import (
	"grovehub/services/identity/internal/app/constants/persona"
	accountdto "grovehub/services/identity/internal/app/dto/account"
)

// Catalog preserves the legacy assignable role contract until tenant roles are
// migrated. A persona is not an authorization grant.
func Catalog() []accountdto.Role {
	definitions := persona.Catalog()
	result := make([]accountdto.Role, 0, len(definitions))
	for _, definition := range definitions {
		result = append(result, accountdto.Role{ID: string(definition.Code), Name: definition.Name})
	}
	return result
}

func IsAssignable(id string) bool {
	return persona.Known(persona.Code(id))
}
