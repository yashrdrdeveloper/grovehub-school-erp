package impl

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/accountrole"
	"grovehub/services/identity/internal/app/repository/accountrole/constants"
)

type Repository struct {
	db *gorm.DB
}

var _ accountrole.Repository = (*Repository)(nil)

func New(db *gorm.DB) accountrole.Repository {
	return &Repository{db: db}
}

type assignmentRow struct {
	AccountID string
	Status entity.AccountStatus
	Version int64
	PrimaryRoleID string
	SecondaryJSON string
}

type candidateRow struct {
	RoleID string
	Status entity.RoleStatus
	Version int64
	CodesJSON string
}

func (repository *Repository) Find(ctx context.Context, tenantID string, accountID string) (accountrole.Assignment, error) {
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID {
		return accountrole.Assignment{}, accountrole.ErrScopeRequired
	}
	id, err := uuid.Parse(accountID)
	if err != nil || id == uuid.Nil {
		return accountrole.Assignment{}, accountrole.ErrNotFound
	}
	if repository.db == nil {
		return accountrole.Assignment{}, transaction.ErrMissingDatabase
	}
	var row assignmentRow
	result := transaction.Database(ctx, repository.db).Raw(constants.FindAssignment, tenantID, id.String()).Scan(&row)
	if result.Error != nil {
		return accountrole.Assignment{}, result.Error
	}
	if result.RowsAffected != 1 {
		return accountrole.Assignment{}, accountrole.ErrNotFound
	}
	assignment := accountrole.Assignment{AccountID: row.AccountID, Status: row.Status, Version: row.Version, PrimaryRoleID: row.PrimaryRoleID}
	if err := json.Unmarshal([]byte(row.SecondaryJSON), &assignment.SecondaryRoleIDs); err != nil {
		return accountrole.Assignment{}, err
	}
	return assignment, nil
}

func (repository *Repository) Candidates(ctx context.Context, tenantID string, roleIDs []string) ([]accountrole.Candidate, error) {
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID {
		return nil, accountrole.ErrScopeRequired
	}
	if repository.db == nil {
		return nil, transaction.ErrMissingDatabase
	}
	if len(roleIDs) == 0 {
		return []accountrole.Candidate{}, nil
	}
	var rows []candidateRow
	result := transaction.Database(ctx, repository.db).Raw(constants.FindCandidates, tenantID, roleIDs).Scan(&rows)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(rows) != len(roleIDs) {
		return nil, accountrole.ErrNotFound
	}
	candidates := make([]accountrole.Candidate, 0, len(rows))
	for _, row := range rows {
		candidate := accountrole.Candidate{RoleID: row.RoleID, Status: row.Status, Version: row.Version}
		if err := json.Unmarshal([]byte(row.CodesJSON), &candidate.Codes); err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate)
	}
	return candidates, nil
}
