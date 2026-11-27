package registry

import (
	"context"
	"errors"

	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
)

var (
	ErrConflict       = errors.New("registered capability has a conflicting module or feature binding")
	ErrFeatureMissing = errors.New("permission feature is not registered")
)

// Batch contains validated, code-owned definitions. No lifecycle or grant state
// is supplied by the caller.
type Batch struct {
	Features    []feature.Definition
	Permissions []permission.Definition
}

type Result struct {
	FeaturesCreated    int
	PermissionsCreated int
}

type Repository interface {
	Register(ctx context.Context, batch Batch) (Result, error)
}
