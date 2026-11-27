// Package delegation evaluates a proposed permission replacement. A successful
// preview is not approval or authorization to persist a change.
package delegation

import (
	"errors"
	"slices"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/services/authorization"
)

var ErrInvalidSet = errors.New("invalid permission set")

type Change struct {
	Added   []permission.Code
	Removed []permission.Code
}

// Preview requires authority over the entire old and proposed set, including
// removals. It uses one effective actor snapshot supplied by the owning service.
func Preview(effective []permission.Code, current []permission.Code, requested []permission.Code) (Change, error) {
	held := make(map[permission.Code]bool, len(effective))
	for _, code := range effective {
		held[code] = true
	}
	if !held[permission.RolePermissionManage] {
		return Change{}, authorization.ErrForbidden
	}
	old, err := permissionSet(current)
	if err != nil {
		return Change{}, err
	}
	next, err := permissionSet(requested)
	if err != nil {
		return Change{}, err
	}
	change := Change{Added: []permission.Code{}, Removed: []permission.Code{}}
	for code := range old {
		if !held[code] || !permission.SchoolDelegable(code) {
			return Change{}, authorization.ErrForbidden
		}
		if !next[code] {
			change.Removed = append(change.Removed, code)
		}
	}
	for code := range next {
		if !held[code] || !permission.SchoolDelegable(code) {
			return Change{}, authorization.ErrForbidden
		}
		if !old[code] {
			change.Added = append(change.Added, code)
		}
	}
	slices.Sort(change.Added)
	slices.Sort(change.Removed)
	return change, nil
}

func permissionSet(codes []permission.Code) (map[permission.Code]bool, error) {
	if len(codes) > len(permission.Catalog()) {
		return nil, ErrInvalidSet
	}
	set := make(map[permission.Code]bool, len(codes))
	for _, code := range codes {
		if _, known := permission.Find(code); !known || set[code] {
			return nil, ErrInvalidSet
		}
		set[code] = true
	}
	return set, nil
}
