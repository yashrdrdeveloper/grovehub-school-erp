package authentication

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
