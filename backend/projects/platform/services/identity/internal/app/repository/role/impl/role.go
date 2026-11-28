package impl

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	coreid "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/id"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/repository/role/constants"
	"grovehub/services/identity/internal/app/repository/role/helpers"
)

type Repository struct {
	db *gorm.DB
}

type detailRow struct {
	entity.Role `gorm:"embedded"`
	PersonaJSON string
}

var _ role.Repository = (*Repository)(nil)

func New(db *gorm.DB) role.Repository {
	return &Repository{db: db}
}

func (repository *Repository) CreateDraft(ctx context.Context, draft role.Draft) (role.Detail, error) {
	draft, err := helpers.NormalizeDraft(draft)
	if err != nil {
		return role.Detail{}, err
	}
	if repository.db == nil {
		return role.Detail{}, coredatabase.ErrMissingDatabase
	}
	identifier, err := coreid.NewUUIDv7()
	if err != nil {
		return role.Detail{}, err
	}
	var detail role.Detail
	err = coredatabase.RunTransaction(ctx, repository.db, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, repository.db)
		normalizedName := strings.ToLower(draft.Name)
		searchDetails := strings.ToLower(strings.TrimSpace(draft.Name + " " + draft.Description))
		if err := db.Exec(constants.InsertDraft, identifier, draft.TenantID, draft.Name, normalizedName, draft.Description, searchDetails, draft.ActorID, draft.ActorID).Error; err != nil {
			return helpers.MapWriteError(err)
		}
		for _, code := range draft.Personas {
			if err := db.Exec(constants.InsertPersona, draft.TenantID, identifier, code).Error; err != nil {
				return err
			}
		}
		var err error
		detail, err = repository.FindByID(tx, draft.TenantID, identifier)
		return err
	})
	if err != nil {
		return role.Detail{}, err
	}
	return detail, nil
}

func (repository *Repository) FindByID(ctx context.Context, tenantID string, roleID string) (role.Detail, error) {
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID {
		return role.Detail{}, role.ErrScopeRequired
	}
	identifier, err := uuid.Parse(roleID)
	if err != nil || identifier == uuid.Nil {
		return role.Detail{}, role.ErrNotFound
	}
	if repository.db == nil {
		return role.Detail{}, coredatabase.ErrMissingDatabase
	}
	var row detailRow
	result := coredatabase.TransactionDB(ctx, repository.db).Raw(constants.FindDetail, tenantID, identifier.String()).Scan(&row)
	if result.Error != nil {
		return role.Detail{}, result.Error
	}
	if result.RowsAffected != 1 {
		return role.Detail{}, role.ErrNotFound
	}
	detail := role.Detail{Role: row.Role}
	if err := json.Unmarshal([]byte(row.PersonaJSON), &detail.Personas); err != nil {
		return role.Detail{}, err
	}
	return detail, nil
}
