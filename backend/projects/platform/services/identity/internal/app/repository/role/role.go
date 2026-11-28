package role

import (
	"context"
	"errors"

	coresearch "grovehub/libs/core/search"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/constants/persona"
	"grovehub/services/identity/internal/app/entity"
)

var (
	ErrScopeRequired   = errors.New("role requires tenant scope")
	ErrInvalidDraft    = errors.New("invalid role draft")
	ErrNotFound        = errors.New("role not found")
	ErrConflict        = errors.New("role name already exists")
	ErrVersionConflict = errors.New("role version conflict")
)

// Draft cannot carry an active status, permissions or account assignments.
type Draft struct {
	TenantID    string
	ActorID     string
	Name        string
	Description string
	Personas    []persona.Code
}

// MetadataUpdate replaces metadata only; grants and lifecycle are untouched.
type MetadataUpdate struct {
	Draft
	RoleID  string
	Version int64
}

type PermissionSet struct {
	RoleID  string
	Version int64
	Codes   []permission.Code
}

type Detail struct {
	Role     entity.Role
	Personas []persona.Code
}

// Repository is internal persistence, not an authorization boundary.
type Repository interface {
	CreateDraft(ctx context.Context, draft Draft) (Detail, error)
	FindByID(ctx context.Context, tenantID string, roleID string) (Detail, error)
	UpdateMetadata(ctx context.Context, update MetadataUpdate) (Detail, error)
	Search(ctx context.Context, tenantID string, query coresearch.Query[string]) (coresearch.QueryResult[Detail], error)
	Permissions(ctx context.Context, tenantID string, roleID string) (PermissionSet, error)
}
