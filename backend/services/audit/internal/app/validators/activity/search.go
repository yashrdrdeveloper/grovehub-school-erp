package activity

import (
	"strings"
	"github.com/google/uuid"
	dto "grovehub/services/audit/internal/app/dto/activity"
	"grovehub/services/audit/internal/app/services/activity"
)

func Normalize(query dto.Search) (dto.Search, error) {
	if query.Limit == 0 { query.Limit = 25 }
	if query.Limit < 1 || query.Limit > 100 || (query.Scope != "" && query.Scope != "ME") || (query.From != nil && query.To != nil && !query.From.Before(*query.To)) { return dto.Search{}, activity.ErrInvalidSearch }
	for _, id := range []string{query.AssetID,query.ActorID,query.PersonID} {
		if id == "" { continue }
		parsed, err := uuid.Parse(id)
		if err != nil || parsed == uuid.Nil || parsed.String() != id { return dto.Search{}, activity.ErrInvalidSearch }
	}
	for _, value := range []string{query.AssetType,query.AuditType,query.Operation,query.Source} {
		if len(value) > 128 || strings.TrimSpace(value) != value { return dto.Search{}, activity.ErrInvalidSearch }
	}
	if query.Cursor != nil {
		id, err := uuid.Parse(query.Cursor.AuditID)
		if err != nil || id == uuid.Nil || id.String() != query.Cursor.AuditID || query.Cursor.OccurredAt.IsZero() { return dto.Search{}, activity.ErrInvalidSearch }
	}
	return query,nil
}
