package operator_test

import (
	"testing"

	dto "grovehub/services/audit/internal/app/dto/operator"
	validator "grovehub/services/audit/internal/app/validators/operator"
)

func TestReceiptFilters(t *testing.T) {
	input := dto.ReceiptSearchInput{SearchInput: dto.SearchInput{PageNumber: 1, RowSize: 25}}
	if !validator.ValidReceiptSearch(input) {
		t.Fatal("global search denied")
	}

	input.SchoolID = "e93696cc-15f0-4e06-9577-0d5cfc2a83c8"
	for _, operation := range []string{"", "RENAME", "ACTIVATE", "DISABLE"} {
		input.Operation = operation
		if !validator.ValidReceiptSearch(input) {
			t.Fatal(operation)
		}
	}

	input.Operation = "DELETE"
	if validator.ValidReceiptSearch(input) {
		t.Fatal("unknown operation")
	}

	input.Operation = ""
	for _, id := range []string{"00000000-0000-0000-0000-000000000000", "E93696CC-15F0-4E06-9577-0D5CFC2A83C8"} {
		input.SchoolID = id
		if validator.ValidReceiptSearch(input) {
			t.Fatal("invalid UUID")
		}
	}

	input.SchoolID = ""
	input.PageNumber = 0
	if validator.ValidReceiptSearch(input) {
		t.Fatal("invalid pagination")
	}
}
