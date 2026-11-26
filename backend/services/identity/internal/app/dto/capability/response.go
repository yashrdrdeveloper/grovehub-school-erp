package capability

import (
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
)

type Discovery struct {
	Permissions []permission.Code `json:"permissions"`
	Actions map[string]bool `json:"actions"`
	ExecutionConfigured bool `json:"executionConfigured"`
	ExecutionRequiresReauthentication bool `json:"executionRequiresReauthentication"`
}

type Permission struct {
	Code permission.Code `json:"code"`
	Name string `json:"name"`
	Feature feature.Code `json:"feature"`
	CanDelegate bool `json:"canDelegate"`
}

type Catalog struct {
	Permissions []Permission `json:"permissions"`
}
