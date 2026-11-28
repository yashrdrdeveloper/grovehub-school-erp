// Package authentication is Identity's public in-process verification contract.
package authentication

import "grovehub/libs/platform/session"

// Client checks opaque session secrets without exposing repositories.
type Client = session.AccessVerifier

// Verifier preserves the existing authentication boundary.
type Verifier = Client
