package impl

import (
	"context"
	"strings"

	"github.com/google/uuid"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/repository/role/constants"
	"grovehub/services/identity/internal/app/repository/role/helpers"
)

func (repository *Repository) UpdateMetadata(ctx context.Context, update role.MetadataUpdate) (role.Detail, error) {
	draft, err := helpers.NormalizeDraft(update.Draft)
	if err != nil {
		return role.Detail{}, err
	}
	identifier, err := uuid.Parse(update.RoleID)
	if err != nil || identifier == uuid.Nil {
		return role.Detail{}, role.ErrNotFound
	}
	if update.Version < 1 {
		return role.Detail{}, role.ErrVersionConflict
	}
	if repository.db == nil {
		return role.Detail{}, coredatabase.ErrMissingDatabase
	}
	var detail role.Detail
	err = coredatabase.RunTransaction(ctx, repository.db, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, repository.db)
		normalizedName := strings.ToLower(draft.Name)
		searchDetails := strings.ToLower(strings.TrimSpace(draft.Name + " " + draft.Description))
		result := db.Exec(constants.UpdateMetadata, draft.Name, normalizedName, draft.Description, searchDetails, draft.ActorID, draft.TenantID, identifier.String(), update.Version)
		if result.Error != nil {
			return helpers.MapWriteError(result.Error)
		}
		if result.RowsAffected != 1 {
			if _, err := repository.FindByID(tx, draft.TenantID, identifier.String()); err != nil {
				return err
			}
			return role.ErrVersionConflict
		}
		if err := db.Exec(constants.DeletePersonas, draft.TenantID, identifier.String()).Error; err != nil {
			return err
		}
		for _, code := range draft.Personas {
			if err := db.Exec(constants.InsertPersona, draft.TenantID, identifier.String(), code).Error; err != nil {
				return err
			}
		}
		var err error
		detail, err = repository.FindByID(tx, draft.TenantID, identifier.String())
		return err
	})
	if err != nil {
		return role.Detail{}, err
	}
	return detail, nil
}
