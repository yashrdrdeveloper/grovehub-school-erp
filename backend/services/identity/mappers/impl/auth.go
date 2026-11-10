package impl

import (
    identitymappers "grovehub/services/identity/mappers"
)

type AuthMapper struct{}

func NewAuthMapper() identitymappers.AuthMapper {
    return &AuthMapper{}
}
