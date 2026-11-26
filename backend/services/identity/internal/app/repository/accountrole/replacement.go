package accountrole

import "context"

type RoleVersion struct {
	RoleID string
	Version int64
}

type Replacement struct {
	TenantID string
	ActorID string
	AccountID string
	Version int64
	PrimaryRoleID string
	SecondaryRoleIDs []string
	RoleVersions []RoleVersion
}

// Writer is an internal persistence primitive. It must not be exposed through
// a service or route before approved execution and legacy cutover are ready.
type Writer interface {
	Replace(ctx context.Context, request Replacement) (Assignment, error)
}
