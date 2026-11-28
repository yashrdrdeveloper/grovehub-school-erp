// Package authentication is Identity's public in-process verification contract.
package authentication

import schoolgateway "grovehub/launcher/services/gateway/session"

// Verifier checks opaque session secrets without exposing repositories.
type Verifier = schoolgateway.AccessVerifier
