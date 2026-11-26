package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"
	"strings"

	"grovehub/services/identity/internal/app/constants/permission"
	accountdto "grovehub/services/identity/internal/app/dto/accountrole"
	dto "grovehub/services/identity/internal/app/dto/approval"
	repo "grovehub/services/identity/internal/app/repository/approval"
	accountvalidator "grovehub/services/identity/internal/app/validators/accountrole"
)

func NormalizePermissions(intent dto.PermissionIntent) (dto.PermissionIntent, error) {
	id, err := NormalizeID(intent.RoleID)
	if err != nil || intent.Version < 1 || intent.Codes == nil || len(intent.Codes) > len(permission.Catalog()) {
		return dto.PermissionIntent{}, repo.ErrInvalid
	}
	codes := append([]permission.Code{}, intent.Codes...)
	slices.Sort(codes)
	for index, code := range codes {
		if !permission.SchoolDelegable(code) || (index > 0 && codes[index-1] == code) {
			return dto.PermissionIntent{}, repo.ErrInvalid
		}
	}
	intent.RoleID = id
	intent.Codes = codes
	return intent, nil
}

func NormalizeAssignment(intent dto.AssignmentIntent) (dto.AssignmentIntent, error) {
	id, err := NormalizeID(intent.AccountID)
	if err != nil || len(intent.RoleVersions) < 1 || len(intent.RoleVersions) > 2*(1+accountvalidator.MaxSecondaryRoles) {
		return dto.AssignmentIntent{}, repo.ErrInvalid
	}
	replacement, err := accountvalidator.NormalizeReplace(AssignmentReplacement(intent))
	if err != nil {
		return dto.AssignmentIntent{}, err
	}
	versions := append([]accountdto.RoleVersion{}, intent.RoleVersions...)
	seen := make(map[string]bool, len(versions))
	for index, version := range versions {
		roleID, err := NormalizeID(version.RoleID)
		if err != nil || version.Version < 1 || seen[roleID] {
			return dto.AssignmentIntent{}, repo.ErrInvalid
		}
		seen[roleID] = true
		versions[index].RoleID = roleID
	}
	for _, roleID := range append([]string{replacement.PrimaryRoleID}, replacement.SecondaryRoleIDs...) {
		if !seen[roleID] {
			return dto.AssignmentIntent{}, repo.ErrInvalid
		}
	}
	slices.SortFunc(versions, func(a, b accountdto.RoleVersion) int { return strings.Compare(a.RoleID, b.RoleID) })
	intent.AccountID = id
	intent.PrimaryRoleID = replacement.PrimaryRoleID
	intent.SecondaryRoleIDs = replacement.SecondaryRoleIDs
	intent.RoleVersions = versions
	return intent, nil
}

func AssignmentReplacement(intent dto.AssignmentIntent) accountdto.Replace {
	return accountdto.Replace{Version: intent.Version, PrimaryRoleID: intent.PrimaryRoleID, SecondaryRoleIDs: intent.SecondaryRoleIDs}
}

func EncodeIntent(intent any) ([]byte, string, error) {
	payload, err := json.Marshal(intent)
	if err != nil {
		return nil, "", err
	}
	digest := sha256.Sum256(payload)
	return payload, hex.EncodeToString(digest[:]), nil
}

func NormalizeReason(reason string, key string) (string, error) {
	reason = strings.TrimSpace(reason)
	if !validText(reason, 2000) || !validText(key, 128) || strings.TrimSpace(key) != key {
		return "", repo.ErrInvalid
	}
	return reason, nil
}
