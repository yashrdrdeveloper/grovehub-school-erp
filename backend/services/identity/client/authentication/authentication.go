// Package authentication is Identity's public in-process verification contract.
package authentication

import coregateway "grovehub/libs/core/gateway"

// Verifier checks opaque session secrets without exposing repositories.
type Verifier = coregateway.AccessVerifier
