package helpers

import (
	coresearch "grovehub/libs/core/search"
	"grovehub/services/audit/internal/app/entity"
	"time"
)

func NextCursor(query coresearch.Query[string], page coresearch.PageRequest, last entity.Activity) (string, error) {
	payload := coresearch.CursorPayload{Field: query.Sort.Field, Direction: query.Sort.Direction, Value: sortValue(last, query.Sort.Field), CursorId: last.AuditID, PageNumber: coresearch.NextPageNumber(page, query.PageNumber), RowSize: query.Limit}
	if payload.PageNumber == 0 {
		payload.RowSize = 0
	}
	for _, rule := range query.SecondarySorts {
		payload.SecondarySorts = append(payload.SecondarySorts, coresearch.Sort{Field: rule.Field, Direction: rule.Direction})
		payload.SecondaryValues = append(payload.SecondaryValues, sortValue(last, rule.Field))
	}
	return coresearch.EncodeCursor(payload)
}

func sortValue(item entity.Activity, field string) string {
	switch field {
	case "auditId":
		return item.AuditID
	case "ingestedAt":
		return item.IngestedAt.UTC().Format(time.RFC3339Nano)
	case "assetType":
		return item.AssetType
	case "operation":
		return item.Operation
	case "source":
		return item.Source
	case "auditType":
		return item.AuditType
	default:
		return item.OccurredAt.UTC().Format(time.RFC3339Nano)
	}
}
