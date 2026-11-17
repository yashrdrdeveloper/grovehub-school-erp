package request

import corerequest "grovehub/libs/core/request"

// Login contains credentials only; no identity or role claims are accepted from the caller.
type Login struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type LoginRequest struct {
    corerequest.Request
    RequestBody Login `validate:"valid"`
}
