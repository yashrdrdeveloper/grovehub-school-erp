package approval

import (
	"strings"

	"github.com/google/uuid"
	coreutils "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	kind "grovehub/services/identity/internal/app/constants/approval"
	dto "grovehub/services/identity/internal/app/dto/approval"
	repo "grovehub/services/identity/internal/app/repository/approval"
)

func normalizeID(_ coresearch.Operator, value string) (string, error) {
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return "", repo.ErrInvalid
	}
	return id.String(), nil
}

func normalizeKind(_ coresearch.Operator, value string) (string, error) {
	if !kind.Known(kind.RequestType(value)) {
		return "", repo.ErrInvalid
	}
	return value, nil
}

func normalizeStatus(_ coresearch.Operator, value string) (string, error) {
	switch kind.Status(value) {
	case kind.Pending, kind.Approved, kind.Rejected, kind.Cancelled, kind.Expired, kind.Applied, kind.Failed, kind.Stale:
		return value, nil
	default:
		return "", repo.ErrInvalid
	}
}

func normalizeTime(_ coresearch.Operator, value string) (string, error) {
	return coreutils.NormalizeRFC3339UTC(value)
}

func cursorValidator(normalize coresearch.ValueNormalizer) coresearch.CursorValueValidator {
	return func(value string, invalid error) error {
		normalized, err := normalize(coresearch.OperatorEquals, value)
		if err != nil || normalized != value {
			return invalid
		}
		return nil
	}
}

func NormalizeSearch(request coresearch.Request) (coresearch.Query[string], error) {
	if strings.ContainsRune(request.SearchText, 0) {
		return coresearch.Query[string]{}, repo.ErrInvalid
	}
	terms, err := coresearch.NormalizeSearchTerms(request.SearchText, 256, 10, repo.ErrInvalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	page, rows, offset, err := coresearch.ResolvePage(request.PageInfo, repo.ErrInvalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	sorts := map[string]coresearch.CursorValueValidator{
		dto.FieldRequestID: cursorValidator(normalizeID), dto.FieldRequestType: cursorValidator(normalizeKind), dto.FieldStatus: cursorValidator(normalizeStatus),
		dto.FieldCreatedAt: coresearch.ValidateTimeCursor, dto.FieldUpdatedAt: coresearch.ValidateTimeCursor, dto.FieldExpiresAt: coresearch.ValidateTimeCursor,
	}
	ordered, err := coresearch.NormalizeSorts(request.Sort, coresearch.QuerySort{Field: dto.FieldCreatedAt, Direction: coresearch.DirectionDescending}, sorts, 5, repo.ErrInvalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	equality := coresearch.OperatorSet(coresearch.OperatorEquals, coresearch.OperatorNotEquals, coresearch.OperatorIn, coresearch.OperatorNotIn)
	fields := map[string]coresearch.FieldSpec{
		dto.FieldRequestID: {Operators: equality, Normalizer: normalizeID},
		dto.FieldRequestType: {Operators: equality, Normalizer: normalizeKind},
		dto.FieldStatus: {Operators: equality, Normalizer: normalizeStatus},
		dto.FieldCreatedAt: {Operators: coresearch.TimeOperators(), Normalizer: normalizeTime},
		dto.FieldUpdatedAt: {Operators: coresearch.TimeOperators(), Normalizer: normalizeTime},
		dto.FieldExpiresAt: {Operators: coresearch.TimeOperators(), Normalizer: normalizeTime},
	}
	filters, err := coresearch.NormalizeFilters(request.Filters, fields, 20, 50, repo.ErrInvalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	cursor, err := coresearch.NormalizeMultiCursor(request.PageInfo.Cursor, ordered, sorts, cursorValidator(normalizeID), repo.ErrInvalid)
	if err != nil {
		return coresearch.Query[string]{}, err
	}
	if err := coresearch.ValidatePageCursor(request.PageInfo, page, rows, repo.ErrInvalid); err != nil {
		return coresearch.Query[string]{}, err
	}
	if request.PageInfo.PageNumber == 0 {
		page = 0
	}
	return coresearch.Query[string]{SearchTerms: terms, Filters: filters, Sort: ordered[0], SecondarySorts: ordered[1:], Cursor: cursor, Limit: rows, PageNumber: page, Offset: offset}, nil
}
