package helpers

import (
	"time"

	coresearch "grovehub/libs/core/search"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
)

func ToSummary(request entity.ApprovalRequest) dto.Summary {
	return dto.Summary{RequestID: request.RequestID, RequesterID: request.RequesterID, RequestType: request.RequestType, Status: request.Status, CurrentStep: request.CurrentStep, StepCount: request.StepCount, Version: request.Version, CreatedAt: request.CreatedAt, UpdatedAt: request.UpdatedAt, ExpiresAt: request.ExpiresAt}
}

func Hydrate(query coresearch.Query[string], rows coresearch.QueryResult[dto.Summary]) (coresearch.Result[dto.Summary], error) {
	result := coresearch.Result[dto.Summary]{Items: rows.Items, PageInfo: coresearch.NewPageResponse(query, rows)}
	if !rows.HasNext || len(rows.Items) == 0 {
		return result, nil
	}
	last := rows.Items[len(rows.Items)-1]
	payload := coresearch.CursorPayload{Field: query.Sort.Field, Direction: query.Sort.Direction, Value: sortValue(last, query.Sort.Field), CursorId: last.RequestID}
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
		return coresearch.Result[dto.Summary]{}, err
	}
	result.PageInfo.NextCursor = cursor
	return result, nil
}

func sortValue(row dto.Summary, field string) string {
	switch field {
	case dto.FieldRequestID:
		return row.RequestID
	case dto.FieldRequestType:
		return string(row.RequestType)
	case dto.FieldStatus:
		return string(row.Status)
	case dto.FieldUpdatedAt:
		return row.UpdatedAt.UTC().Format(time.RFC3339Nano)
	case dto.FieldExpiresAt:
		return row.ExpiresAt.UTC().Format(time.RFC3339Nano)
	default:
		return row.CreatedAt.UTC().Format(time.RFC3339Nano)
	}
}
