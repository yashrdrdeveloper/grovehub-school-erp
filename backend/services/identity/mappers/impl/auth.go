package impl

import (
    identitymappers "grovehub/services/identity/mappers"
)

type AuthMapper struct{}

func NewAuthMapper() identitymappers.AuthMapper {
    return &AuthMapper{}
}

// nullableString converts optional model storage into its API representation.
func nullableString(value *string) string {
    if value == nil { return "" }
    return *value
}
