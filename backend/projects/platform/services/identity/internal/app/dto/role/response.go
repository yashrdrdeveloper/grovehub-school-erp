package role

import (
	"time"

	"grovehub/services/identity/internal/app/constants/persona"
)

type Detail struct {
	RoleID      string         `json:"roleId"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Personas    []persona.Code `json:"personas"`
	Version     int64          `json:"version"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}
