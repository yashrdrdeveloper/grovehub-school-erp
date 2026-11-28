package workload

import (
	"time"

	"grovehub/libs/platform/scope"
	"grovehub/services/registry/internal/app/constants/status"
)

type CreateInput struct {
	ProjectID *string `json:"projectId,omitempty"`
	Issuer    string  `json:"issuer"`
	Subject   string  `json:"subject"`
}

type StatusInput struct {
	ExpectedVersion int64         `json:"expectedVersion"`
	Status          status.Status `json:"status"`
}

type KeyInput struct {
	KeyID      string    `json:"keyId"`
	PublicKey  string    `json:"publicKey"`
	ValidFrom  time.Time `json:"validFrom"`
	ValidUntil time.Time `json:"validUntil"`
}

type GrantInput struct {
	Scope           scope.Level `json:"scope"`
	TargetProjectID string      `json:"targetProjectId,omitempty"`
	TargetTenantID  string      `json:"targetTenantId,omitempty"`
	Capability      string      `json:"capability"`
	AllowDelegation bool        `json:"allowDelegation"`
}
