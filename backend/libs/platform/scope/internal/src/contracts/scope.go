package contracts

import "errors"

type Level string

const (
	Platform Level = "PLATFORM"
	Project Level = "PROJECT"
	Tenant Level = "TENANT"
)

type PrincipalType string

const (
	User     PrincipalType = "USER"
	Workload PrincipalType = "WORKLOAD"
	System   PrincipalType = "SYSTEM"
)

var (
	ErrInvalid = errors.New("invalid execution scope")
	ErrDenied  = errors.New("execution scope denied")
)

// Target denotes the resource ownership boundary, not the caller's identity.
type Target struct {
	Level     Level
	ProjectID string
	TenantID  string
}

// Execution describes a single verified operation. CallerProjectID identifies
// the credential owner; Target identifies the authorized resource owner.
// System and platform-owned workloads may have no caller project.
type Execution struct {
	Target          Target
	Capability      string
	CallerProjectID string
	PrincipalType   PrincipalType
	PrincipalID     string
	WorkloadID      string
	SessionID       string
	CorrelationID   string
	CausationID     string
}

// Grant must be obtained from trusted, current authorization state.
type Grant struct {
	CallerProjectID string
	WorkloadID      string
	Target          Target
	Capability      string
	AllowDelegation bool
}
