package search

import "testing"

func TestResolvePage(t *testing.T) {
	tests := []struct {
		name string
		request PageRequest
		page, rows, offset int
		invalid bool
	}{
		{name: "default", request: PageRequest{}, page: 1, rows: 25},
		{name: "legacy", request: PageRequest{Rows: 50}, page: 1, rows: 50},
		{name: "direct jump", request: PageRequest{PageNumber: 12, RowSize: 25}, page: 12, rows: 25, offset: 275},
		{name: "cursor takes precedence", request: PageRequest{PageNumber: 12, RowSize: 25, Cursor: "cursor"}, page: 12, rows: 25},
		{name: "conflicting sizes", request: PageRequest{Rows: 50, RowSize: 25}, invalid: true},
		{name: "invalid page", request: PageRequest{PageNumber: -1}, invalid: true},
		{name: "invalid size", request: PageRequest{RowSize: 101}, invalid: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			page, rows, offset, err := ResolvePage(tc.request, ErrInvalidQuery)
			if tc.invalid {
				if err == nil { t.Fatal("expected validation error") }
				return
			}
			if err != nil || page != tc.page || rows != tc.rows || offset != tc.offset {
				t.Fatalf("page=%d rows=%d offset=%d err=%v", page, rows, offset, err)
			}
		})
	}
}

func TestValidatePageCursor(t *testing.T) {
	token, err := EncodeCursor(CursorPayload{PageNumber: 3, RowSize: 25})
	if err != nil { t.Fatal(err) }
	request := PageRequest{PageNumber: 3, RowSize: 25, Cursor: token}
	if err := ValidatePageCursor(request, 3, 25, ErrInvalidQuery); err != nil { t.Fatal(err) }
	if err := ValidatePageCursor(request, 4, 25, ErrInvalidQuery); err == nil { t.Fatal("wrong page accepted") }
	if err := ValidatePageCursor(request, 3, 50, ErrInvalidQuery); err == nil { t.Fatal("wrong row size accepted") }
}

func TestNewPageResponse(t *testing.T) {
	total := int64(101)
	response := NewPageResponse(Query[string]{PageNumber: 3, Limit: 25}, QueryResult[string]{Items: []string{"a"}, TotalRecords: &total, HasNext: true})
	if response.TotalPages == nil || *response.TotalPages != 5 || response.PageNumber != 3 || response.RowSize != 25 || !response.HasPrevious || !response.HasNext {
		t.Fatalf("unexpected response: %+v", response)
	}
}
