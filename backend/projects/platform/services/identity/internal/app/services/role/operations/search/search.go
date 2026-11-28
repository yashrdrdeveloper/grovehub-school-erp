package search

import (
	"time"

	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	roledto "grovehub/services/identity/internal/app/dto/role"
	"grovehub/services/identity/internal/app/entity"
	rolemapper "grovehub/services/identity/internal/app/mappers/role"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
)

func Hydrate(query coresearch.Query[string], rows coresearch.QueryResult[rolerepo.Detail]) (coresearch.Result[roledto.Detail], error) {
	result := coresearch.Result[roledto.Detail]{Items: make([]roledto.Detail, 0, len(rows.Items)), PageInfo: coresearch.NewPageResponse(query, rows)}
	for _, detail := range rows.Items {
		result.Items = append(result.Items, rolemapper.ToDetail(detail))
	}
	if !rows.HasNext || len(rows.Items) == 0 {
		return result, nil
	}
	last := rows.Items[len(rows.Items)-1].Role
	payload := coresearch.CursorPayload{
		Field: query.Sort.Field, Direction: query.Sort.Direction,
		Value: sortValue(last, query.Sort.Field), CursorId: last.RoleID,
	}
	if query.PageNumber > 0 {
		payload.PageNumber = query.PageNumber + 1
		payload.RowSize = query.Limit
	}
	for _, rule := range query.SecondarySorts {
		payload.SecondarySorts = append(payload.SecondarySorts, coresearch.Sort{Field: rule.Field, Direction: rule.Direction})
		payload.SecondaryValues = append(payload.SecondaryValues, sortValue(last, rule.Field))
	}
	cursor, err := coresearch.EncodeCursor(payload)
	if err != nil {
		return coresearch.Result[roledto.Detail]{}, err
	}
	result.PageInfo.NextCursor = cursor
	return result, nil
}

func sortValue(role entity.Role, field string) string {
	switch field {
	case roledto.FieldRoleID:
		return role.RoleID
	case roledto.FieldName:
		return role.NormalizedName
	case roledto.FieldStatus:
		return string(role.Status)
	case roledto.FieldUpdatedAt:
		return role.UpdatedAt.UTC().Format(time.RFC3339Nano)
	default:
		return role.CreatedAt.UTC().Format(time.RFC3339Nano)
	}
}
