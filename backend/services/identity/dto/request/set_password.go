package request

import corerequest "grovehub/libs/core/request"

// The setup token is a one-time bearer secret delivered through a trusted email.
type SetPassword struct {
    Token string `json:"token"`
    Password string `json:"password"`
}

type SetPasswordRequest struct {
    corerequest.Request
    RequestBody SetPassword `validate:"valid"`
}
