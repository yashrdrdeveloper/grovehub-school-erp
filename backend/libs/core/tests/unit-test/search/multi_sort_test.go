package search_test

import (
	"testing"

	. "grovehub/libs/core/search"
)

func TestNormalizeSortsAndMultiCursor(t *testing.T) {
	allowed := map[string]CursorValueValidator{"status": ValidateNonEmptyCursor, "createdAt": ValidateTimeCursor}
	fallback := QuerySort{Field: "createdAt", Direction: DirectionDescending}
	request := []Sort{{Field: "status", Direction: DirectionAscending}, {Field: "createdAt", Direction: DirectionDescending}}
	rules, err := NormalizeSorts(request, fallback, allowed, 4, ErrInvalidQuery)
	if err != nil || len(rules) != 2 || rules[0].Field != "status" {
		t.Fatalf("sort order: %v %v", rules, err)
	}
	if _, err := NormalizeSorts([]Sort{request[0], request[0]}, fallback, allowed, 4, ErrInvalidQuery); err == nil {
		t.Fatal("duplicate sort accepted")
	}
	id := "019ab812-3456-7890-8123-123456789abc"
	token, err := EncodeCursor(CursorPayload{Field: "status", Direction: DirectionAscending, Value: "ACTIVE", CursorId: id, SecondarySorts: request[1:], SecondaryValues: []string{"2026-11-18T00:00:00Z"}})
	if err != nil {
		t.Fatal(err)
	}
	cursor, err := NormalizeMultiCursor(token, rules, allowed, ValidateUUIDv7Cursor, ErrInvalidQuery)
	if err != nil || len(cursor.SecondaryValues) != 1 || cursor.CursorId != id {
		t.Fatalf("cursor: %v %v", cursor, err)
	}
	rules[1].Direction = DirectionAscending
	if _, err := NormalizeMultiCursor(token, rules, allowed, ValidateUUIDv7Cursor, ErrInvalidQuery); err == nil {
		t.Fatal("changed direction accepted")
	}
}
