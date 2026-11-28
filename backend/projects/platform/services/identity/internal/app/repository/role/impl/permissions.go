package impl

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/repository/role/constants"
)

type permissionRow struct {
	RoleID    string
	Version   int64
	CodesJSON string
}

func (repository *Repository) Permissions(ctx context.Context, tenantID string, roleID string) (role.PermissionSet, error) {
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID {
		return role.PermissionSet{}, role.ErrScopeRequired
	}
	id, err := uuid.Parse(roleID)
	if err != nil || id == uuid.Nil {
		return role.PermissionSet{}, role.ErrNotFound
	}
	if repository.db == nil {
		return role.PermissionSet{}, coredatabase.ErrMissingDatabase
	}
	var row permissionRow
	result := coredatabase.TransactionDB(ctx, repository.db).Raw(constants.FindPermissions, tenantID, id.String()).Scan(&row)
	if result.Error != nil {
		return role.PermissionSet{}, result.Error
	}
	if result.RowsAffected != 1 {
		return role.PermissionSet{}, role.ErrNotFound
	}
	set := role.PermissionSet{RoleID: row.RoleID, Version: row.Version}
	if err := json.Unmarshal([]byte(row.CodesJSON), &set.Codes); err != nil {
		return role.PermissionSet{}, err
	}
	return set, nil
}
