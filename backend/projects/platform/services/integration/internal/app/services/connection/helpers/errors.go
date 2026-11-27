package helpers

import "errors"

// Canonical connection service errors are shared by independent operations.
var (
    ErrInvalidConnection = errors.New("invalid integration connection")
    ErrConnectionNotReady = errors.New("provider authorization is not implemented")
    ErrOAuthUnavailable = errors.New("provider OAuth is not configured")
    ErrOAuthInvalid = errors.New("invalid OAuth callback")
    ErrConnectionInactive = errors.New("connection not active")
    ErrPushUnavailable = errors.New("Gmail Pub/Sub push is not configured")
)
