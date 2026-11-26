package role

import "grovehub/services/identity/internal/app/constants/persona"

// CreateDraft intentionally excludes tenant, actor, status and permission grants.
type CreateDraft struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Personas    []persona.Code `json:"personas"`
}
