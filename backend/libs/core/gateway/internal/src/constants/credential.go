package constants

import "time"

const (
	CredentialVersion = "ghw1"
	MaxCredentialSize = 8192
	MaxKeyIDLength = 128
	MaxCredentialLifetime = 2 * time.Minute
	OutboundCredentialLifetime = time.Minute
	HTTPClientTimeout = 15 * time.Second
)
