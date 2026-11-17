package impl

import (
    legacy "grovehub/services/identity/controller/impl"
    protected "grovehub/services/identity/controller/protected"
    identityservices "grovehub/services/identity/services"
    identityvalidators "grovehub/services/identity/validators"
)

// NewAuthController adapts the existing registration REST pipeline while the
// controller tree is reorganized. Authorization is enforced in the service.
func NewAuthController(service identityservices.AuthService, validator identityvalidators.AuthValidator) protected.AuthController {
    return legacy.NewAuthController(service, validator)
}
