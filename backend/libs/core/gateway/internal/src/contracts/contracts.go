// Package contracts defines gateway tiers and workload authorization contracts.
package contracts

import (
	"context"
	"crypto/ed25519"
	"errors"
	"net/http"
	"time"
)

type Tier string

const (
	Public     Tier = "public"
	Protected  Tier = "protected"
	Restricted Tier = "restricted"
)

type WorkloadVerifier interface {
	VerifyWorkload(ctx context.Context, request *http.Request) (context.Context, error)
}

var ErrWorkloadDenied = errors.New("workload authorization denied")

// WorkloadClaims must originate from a cryptographically verified credential,
// never directly from request headers or an unverified token.
type WorkloadClaims struct {
	Subject string
	Audience string
	ExpiresAt time.Time
	Scopes []string
	TenantID string
	ActorGID string
	SessionID string
}

// CredentialVerifier verifies credential authenticity, issuer and algorithm.
// ScopedWorkloadVerifier additionally enforces audience, expiry and route scope.
type CredentialVerifier interface {
	VerifyCredential(context.Context, string) (WorkloadClaims, error)
}

type RoutePolicy struct {
	Scope string
	RequireActor bool
}

type TrustedWorkloadKey struct {
	PublicKey ed25519.PublicKey
	Issuer string
	Subject string
	Scopes []string
	AllowDelegation bool
}

type WorkloadKeyConfig struct {
	KeyID string
	PublicKey string
	Issuer string
	Subject string
	Scopes []string
	AllowDelegation bool
}
