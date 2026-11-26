package audit

import (
	"slices"

	"github.com/google/uuid"
	dto "grovehub/services/identity/internal/app/dto/audit"
	"grovehub/services/identity/internal/app/services/audit"
)

func sourceValid(source string) bool {
	return slices.Contains([]string{"APPROVAL","RESOURCE","ROLE_METADATA","ACCOUNT_METADATA","REGISTRATION"}, source)
}

func validID(value string) bool {
	id, err := uuid.Parse(value)
	return err == nil && id != uuid.Nil && id.String() == value
}

func Normalize(query dto.Search) (dto.Search, error) {
	if query.Limit == 0 { query.Limit = 25 }
	if query.Limit < 1 || query.Limit > 100 ||
		(query.Source != "" && !sourceValid(query.Source)) ||
		(query.ResourceType != "" && !slices.Contains([]string{"ACCOUNT","ROLE","APPROVAL"}, query.ResourceType)) ||
		(query.ResourceID != "" && !validID(query.ResourceID)) ||
		(query.ActorID != "" && !validID(query.ActorID)) ||
		(query.Operation != "" && !slices.Contains([]string{"CREATED","UPDATED","APPROVE","REJECT","CANCELLED","APPLIED","STALE","EXPIRED","ROLE_LIFECYCLE_CHANGE","ROLE_PERMISSION_CHANGE","ROLE_ASSIGNMENT_CHANGE","ACCOUNT_LIFECYCLE_CHANGE"}, query.Operation)) ||
		(query.From != nil && query.To != nil && !query.From.Before(*query.To)) {
		return dto.Search{}, audit.ErrInvalidSearch
	}
	if query.Cursor != nil && (query.Cursor.OccurredAt.IsZero() || !sourceValid(query.Cursor.Source) || !validID(query.Cursor.EventID)) {
		return dto.Search{}, audit.ErrInvalidSearch
	}
	return query, nil
}
