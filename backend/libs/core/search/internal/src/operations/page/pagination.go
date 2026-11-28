package page

import (
	"math"
	"strings"

	"gorm.io/gorm"

	"grovehub/libs/core/search/internal/src/contracts"
	"grovehub/libs/core/search/internal/src/operations/cursor"
	"grovehub/libs/core/search/internal/src/operations/validation"
)

// ResolvePage keeps legacy rows requests valid and uses offsets only without a cursor.
func ResolvePage(request contracts.PageRequest, invalid error) (int, int, int, error) {
	if invalid == nil { invalid = contracts.ErrInvalidQuery }
	if request.PageNumber < 0 || (request.RowSize != 0 && request.Rows != 0 && request.RowSize != request.Rows) { return 0, 0, 0, invalid }
	size := request.RowSize
	if size == 0 { size = request.Rows }
	rows, err := validation.ResolveRows(size, invalid)
	if err != nil { return 0, 0, 0, err }
	page := request.PageNumber
	if page == 0 { page = 1 }
	if page == math.MaxInt || page-1 > math.MaxInt/rows { return 0, 0, 0, invalid }
	offset := 0
	if strings.TrimSpace(request.Cursor) == "" { offset = (page-1)*rows }
	return page, rows, offset, nil
}

// ValidatePageCursor prevents using a numbered cursor for another page or size.
// Legacy unnumbered requests and cursors remain supported for existing clients.
func ValidatePageCursor(request contracts.PageRequest, page int, rows int, invalid error) error {
	if invalid == nil { invalid = contracts.ErrInvalidQuery }
	if strings.TrimSpace(request.Cursor) == "" { return nil }
	payload, err := cursor.DecodeCursor(request.Cursor)
	if err != nil { return invalid }
	if request.PageNumber > 0 && (payload.PageNumber != page || payload.RowSize != rows) { return invalid }
	return nil
}

func NextPageNumber(request contracts.PageRequest, page int) int {
	if request.PageNumber > 0 { return page+1 }
	return 0
}

func ApplyPage[TId any](db *gorm.DB, query contracts.Query[TId]) *gorm.DB {
	if query.Cursor == nil && query.Offset > 0 { db = db.Offset(query.Offset) }
	return db.Limit(query.Limit+1)
}

func NewPageResponse[TId any, T any](query contracts.Query[TId], result contracts.QueryResult[T]) contracts.PageResponse {
	response := contracts.PageResponse{Rows: len(result.Items), TotalRecords: result.TotalRecords, HasNext: result.HasNext}
	if query.PageNumber == 0 { return response }
	response.PageNumber = query.PageNumber
	response.RowSize = query.Limit
	response.HasPrevious = query.PageNumber > 1
	if result.TotalRecords != nil && query.Limit > 0 {
		pages := *result.TotalRecords / int64(query.Limit)
		if *result.TotalRecords % int64(query.Limit) != 0 { pages++ }
		response.TotalPages = &pages
	}
	return response
}
