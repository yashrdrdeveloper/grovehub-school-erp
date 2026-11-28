package impl

import (
	"context"
	"encoding/json"
	"strings"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	coresearch "grovehub/libs/core/search"
	"grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/repository/role/constants"
	searchconfig "grovehub/services/identity/internal/app/repository/role/helpers/search"
)

func (repository *Repository) Search(ctx context.Context, tenantID string, query coresearch.Query[string]) (coresearch.QueryResult[role.Detail], error) {
	if strings.TrimSpace(tenantID) == "" || strings.TrimSpace(tenantID) != tenantID {
		return coresearch.QueryResult[role.Detail]{}, role.ErrScopeRequired
	}
	if repository.db == nil {
		return coresearch.QueryResult[role.Detail]{}, coredatabase.ErrMissingDatabase
	}
	db := coredatabase.TransactionDB(ctx, repository.db).Table("identity.roles AS r").Where("r.tenant_id = ?", tenantID)
	// Legacy cursor requests also receive the filtered total, as Account Search does.
	var legacyTotal *int64
	if query.PageNumber == 0 {
		countQuery := query
		countQuery.Cursor = nil
		countQuery.Offset = 0
		countDB, err := coresearch.Apply(db.Session(&gorm.Session{}), countQuery, searchconfig.RoleConfig)
		if err != nil {
			return coresearch.QueryResult[role.Detail]{}, err
		}
		var count int64
		if err := countDB.Select("COUNT(*)").Limit(-1).Offset(-1).Count(&count).Error; err != nil {
			return coresearch.QueryResult[role.Detail]{}, err
		}
		legacyTotal = &count
	}
	rows, err := coresearch.Execute[string, detailRow](db.Select(constants.DetailProjection), query, searchconfig.RoleConfig)
	if err != nil {
		return coresearch.QueryResult[role.Detail]{}, err
	}
	if legacyTotal != nil {
		rows.TotalRecords = legacyTotal
	}
	result := coresearch.QueryResult[role.Detail]{Items: make([]role.Detail, 0, len(rows.Items)), HasNext: rows.HasNext, TotalRecords: rows.TotalRecords}
	for _, row := range rows.Items {
		detail := role.Detail{Role: row.Role}
		if err := json.Unmarshal([]byte(row.PersonaJSON), &detail.Personas); err != nil {
			return coresearch.QueryResult[role.Detail]{}, err
		}
		result.Items = append(result.Items, detail)
	}
	return result, nil
}
