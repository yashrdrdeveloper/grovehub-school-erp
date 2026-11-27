package helpers

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"grovehub/services/identity/internal/app/constants/persona"
	"grovehub/services/identity/internal/app/repository/role"
)

func NormalizeDraft(draft role.Draft) (role.Draft, error) {
	if strings.TrimSpace(draft.TenantID) == "" || strings.TrimSpace(draft.TenantID) != draft.TenantID {
		return role.Draft{}, role.ErrScopeRequired
	}
	actorID, err := uuid.Parse(draft.ActorID)
	if err != nil || actorID == uuid.Nil {
		return role.Draft{}, role.ErrInvalidDraft
	}
	draft.Name = strings.TrimSpace(draft.Name)
	draft.Description = strings.TrimSpace(draft.Description)
	if draft.Name == "" || len(draft.Personas) == 0 || len(draft.Personas) > len(persona.Catalog()) {
		return role.Draft{}, role.ErrInvalidDraft
	}
	seen := make(map[persona.Code]bool, len(draft.Personas))
	for _, code := range draft.Personas {
		if !persona.Known(code) || seen[code] {
			return role.Draft{}, role.ErrInvalidDraft
		}
		seen[code] = true
	}
	return draft, nil
}

func MapWriteError(err error) error {
	var postgres *pgconn.PgError
	if errors.As(err, &postgres) && postgres.Code == "23505" && postgres.ConstraintName == "roles_tenant_id_normalized_name_key" {
		return role.ErrConflict
	}
	return err
}
