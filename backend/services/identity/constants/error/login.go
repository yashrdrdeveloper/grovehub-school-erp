package errorcode

import (
    "net/http"
    coreerrors "grovehub/libs/core/errors"
)

var (
    InvalidLogin = coreerrors.NewDefinition("IDENTITY_LOGIN_INVALID", "Login request is invalid", http.StatusBadRequest)
    InvalidCredentials = coreerrors.NewDefinition("IDENTITY_LOGIN_INVALID_CREDENTIALS", "Invalid email or password", http.StatusUnauthorized)
    AccountNotActive = coreerrors.NewDefinition("IDENTITY_LOGIN_ACCOUNT_NOT_ACTIVE", "Account is not active", http.StatusForbidden)
    AuthenticationUnavailable = coreerrors.NewDefinition("IDENTITY_AUTHENTICATION_UNAVAILABLE", "Authentication is not configured", http.StatusServiceUnavailable)
    RegistrationUnauthorized = coreerrors.NewDefinition("IDENTITY_REGISTRATION_FORBIDDEN", "Not authorized to register accounts", http.StatusForbidden)
    RegistrationUnavailable = coreerrors.NewDefinition("IDENTITY_REGISTRATION_UNAVAILABLE", "Registration invitation delivery is not configured", http.StatusServiceUnavailable)
)
