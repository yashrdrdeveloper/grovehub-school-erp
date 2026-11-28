// Package authorization exposes the Registry's live workload verification boundary.
package authorization

import (
	"context"
	"time"

	"grovehub/libs/platform/scope"
)

// Request is constructed by a trusted provider-owned route policy. Target and
// capability must never be selected solely by an unverified token or header.
type Request struct {
	Token               string
	Target              scope.Target
	Capability          string
	Method              string
	RequestURI          string
	AllowInactiveTarget bool
}

// Client verifies nondelegated workload credentials for one exact operation.
type Client interface {
	Verify(ctx context.Context, request Request) (context.Context, error)
	TenantState(ctx context.Context, projectID, tenantID string) (TenantState, error)
}

type TenantState struct {
	Active            bool
	SessionGeneration int64
}

// Credential is an internal lookup result. It contains only public material.
type Credential struct {
	WorkloadID        string
	CallerProjectID   string
	Issuer            string
	Subject           string
	PublicKey         []byte
	ValidFrom         time.Time
	ValidUntil        time.Time
	WorkloadUpdatedAt time.Time
}
