// Package gateway exposes application-neutral route tiers and workload verification.
package gateway

import (
	"context"
	"crypto/ed25519"
	"net/http"

	"grovehub/libs/core/gateway/internal/src/config"
	"grovehub/libs/core/gateway/internal/src/contracts"
	"grovehub/libs/core/gateway/internal/src/operations/credential"
	"grovehub/libs/core/gateway/internal/src/operations/restricted"
	"grovehub/libs/core/gateway/internal/src/operations/transport"
	"grovehub/libs/core/gateway/internal/src/operations/workload"
)

type Tier = contracts.Tier

const (
	Public     = contracts.Public
	Protected  = contracts.Protected
	Restricted = contracts.Restricted
)

type WorkloadVerifier = contracts.WorkloadVerifier
type WorkloadClaims = contracts.WorkloadClaims
type CredentialVerifier = contracts.CredentialVerifier
type RoutePolicy = contracts.RoutePolicy
type TrustedWorkloadKey = contracts.TrustedWorkloadKey
type WorkloadKeyConfig = contracts.WorkloadKeyConfig
type Ed25519Credentials = credential.Ed25519Credentials
type ScopedWorkloadVerifier = workload.ScopedWorkloadVerifier
type WorkloadTransport = transport.WorkloadTransport

var ErrWorkloadDenied = contracts.ErrWorkloadDenied

func NewEd25519Credentials(keys []WorkloadKeyConfig) (Ed25519Credentials, error) {
	return config.NewEd25519Credentials(keys)
}

func SignWorkload(keyID, issuer string, key ed25519.PrivateKey, claims WorkloadClaims) (string, error) {
	return credential.SignWorkload(keyID, issuer, key, claims)
}

func RestrictedOnly(mux *http.ServeMux, downstream http.Handler) http.Handler {
	return restricted.RestrictedOnly(mux, downstream)
}

func NewWorkloadHTTPClient(value *WorkloadTransport) (*http.Client, error) {
	return transport.NewWorkloadHTTPClient(value)
}

func WorkloadFrom(ctx context.Context) (WorkloadClaims, bool) {
	return workload.WorkloadFrom(ctx)
}

var (
	_ WorkloadVerifier = ScopedWorkloadVerifier{}
	_ CredentialVerifier = Ed25519Credentials{}
	_ http.RoundTripper = (*WorkloadTransport)(nil)
)
