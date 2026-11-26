package helpers

import (
	"encoding/hex"
	"math"
	"strings"

	"github.com/google/uuid"

	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/repository/approvalpolicy"
)

func canonicalID(value string) bool {
	id, err := uuid.Parse(value)
	return err == nil && id != uuid.Nil && id.String() == value
}

func Validate(manifest approvalpolicy.Manifest, operatorID string, digest string) error {
	_, err := hex.DecodeString(digest)
	if err != nil || len(digest) != 64 || strings.ToLower(digest) != digest || !canonicalID(operatorID) || !canonicalID(manifest.AuthorID) || manifest.TenantID == "" || strings.TrimSpace(manifest.TenantID) != manifest.TenantID || strings.ContainsRune(manifest.TenantID, 0) {
		return approvalpolicy.ErrInvalid
	}
	if len(manifest.Policies) != 4 {
		return approvalpolicy.ErrInvalid
	}
	seen := make(map[kind.RequestType]bool, 4)
	for _, revision := range manifest.Policies {
		switch revision.Type {
		case kind.AccountLifecycleChange, kind.RoleLifecycleChange, kind.RolePermissionChange, kind.RoleAssignmentChange:
		default:
			return approvalpolicy.ErrInvalid
		}
		if seen[revision.Type] || revision.ExpectedVersion < 0 || revision.ExpectedVersion == math.MaxInt64 || len(revision.Steps) < 1 || len(revision.Steps) > kind.MaxSteps {
			return approvalpolicy.ErrInvalid
		}
		seen[revision.Type] = true
		for _, step := range revision.Steps {
			if !permission.SchoolDelegable(step.Permission) || (step.RoleID != nil && !canonicalID(*step.RoleID)) {
				return approvalpolicy.ErrInvalid
			}
		}
	}
	return nil
}
