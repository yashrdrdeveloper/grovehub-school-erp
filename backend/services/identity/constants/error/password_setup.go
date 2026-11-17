package errorcode

import (
    "net/http"
    coreerrors "grovehub/libs/core/errors"
)

var (
    InvalidPasswordSetup = coreerrors.NewDefinition("IDENTITY_PASSWORD_SETUP_INVALID", "Password setup request is invalid", http.StatusBadRequest)
    InvalidPasswordSetupToken = coreerrors.NewDefinition("IDENTITY_PASSWORD_SETUP_TOKEN_INVALID", "Password setup link is invalid or expired", http.StatusBadRequest)
    PasswordSetupFailed = coreerrors.NewDefinition("IDENTITY_PASSWORD_SETUP_FAILED", "Unable to complete password setup", http.StatusInternalServerError)
)
