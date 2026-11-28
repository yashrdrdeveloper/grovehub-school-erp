package activity

import (
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	dto "grovehub/services/audit/internal/app/dto/activity"
	"grovehub/services/audit/internal/app/services/activity"
	"grovehub/services/audit/internal/app/validators/activity/constants"
	"grovehub/services/audit/internal/app/validators/activity/helpers"
)

func sortFields() map[string]coresearch.CursorValueValidator {
	return map[string]coresearch.CursorValueValidator{
		"occurredAt": coresearch.ValidateTimeCursor,
		"ingestedAt": coresearch.ValidateTimeCursor,
		"auditId":    helpers.ValidateID,
		"assetType":  coresearch.ValidateNonEmptyCursor,
		"operation":  coresearch.ValidateNonEmptyCursor,
		"source":     coresearch.ValidateNonEmptyCursor,
		"auditType":  coresearch.ValidateNonEmptyCursor,
	}
}

func filterFields() map[string]coresearch.FieldSpec {
	fields := map[string]coresearch.FieldSpec{}

	for _, field := range []string{"auditId", "assetId", "actorId", "personId", "subjectAccountId", "occurrenceId"} {
		fields[field] = coresearch.FieldSpec{Operators: coresearch.EqualityOperators(), Normalizer: helpers.NormalizeID}
	}
	for _, field := range []string{"assetType", "auditType", "source", "operation"} {
		fields[field] = coresearch.FieldSpec{Operators: coresearch.TextOperators(), Normalizer: helpers.NormalizeText}
	}
	for _, field := range []string{"occurredAt", "ingestedAt"} {
		fields[field] = coresearch.FieldSpec{Operators: coresearch.TimeOperators(), Normalizer: helpers.NormalizeTime}
	}

	return fields
}

func Normalize(input dto.Search) (coresearch.Query[string], error) {
	invalid := activity.ErrInvalidSearch
	if input.Scope != "" && input.Scope != constants.ScopeMe {
		return coresearch.Query[string]{}, invalid
	}

	terms, err := coresearch.NormalizeSearchTerms(input.SearchText, constants.MaxSearchTextBytes, constants.MaxSearchTerms, invalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}

	page, rows, offset, err := coresearch.ResolvePage(input.PageInfo, invalid)
	if err != nil || offset > constants.MaxOffset {
		return coresearch.Query[string]{}, invalid
	}

	fields := sortFields()
	sorts, err := coresearch.NormalizeSorts(input.Sort, coresearch.QuerySort{Field: "occurredAt", Direction: coresearch.DirectionDescending}, fields, constants.MaxSorts, invalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}

	filters, err := coresearch.NormalizeFilters(input.Filters, filterFields(), constants.MaxFilters, constants.MaxFilterValues, invalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}

	cursor, err := coresearch.NormalizeMultiCursor(input.PageInfo.Cursor, sorts, fields, helpers.ValidateID, invalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	if err := coresearch.ValidatePageCursor(input.PageInfo, page, rows, invalid); err != nil {
		return coresearch.Query[string]{}, err
	}

	if input.PageInfo.PageNumber == 0 {
		page = 0
	}

	return coresearch.Query[string]{
		SearchTerms: terms, Filters: filters, Sort: sorts[0], SecondarySorts: sorts[1:],
		Cursor: cursor, Limit: rows, PageNumber: page, Offset: offset,
	}, nil
}
