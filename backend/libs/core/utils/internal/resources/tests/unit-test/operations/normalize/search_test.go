package normalize_test

import (
	"errors"
	"testing"

	"grovehub/libs/core/utils"
)

func TestStrictNormalization(t *testing.T) {
	const id = "e93696cc-15f0-4e06-9577-0d5cfc2a83c8"
	if err := utils.ValidateUUID(id); err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{"", "00000000-0000-0000-0000-000000000000", "E93696CC-15F0-4E06-9577-0D5CFC2A83C8", " " + id} {
		if !errors.Is(utils.ValidateUUID(value), utils.ErrInvalidUUID) {
			t.Fatalf("accepted invalid UUID: %q", value)
		}
	}
	if value, err := utils.NormalizeText(" text ", 4); err != nil || value != "text" {
		t.Fatalf("text: %q, %v", value, err)
	}
	for _, value := range []string{" ", "too long"} {
		if _, err := utils.NormalizeText(value, 4); !errors.Is(err, utils.ErrInvalidText) {
			t.Fatalf("accepted invalid text: %q", value)
		}
	}
	if value, err := utils.NormalizeStrictRFC3339UTC("2026-11-29T05:30:00.123456789+05:30"); err != nil || value != "2026-11-29T00:00:00.123456789Z" {
		t.Fatalf("timestamp: %q, %v", value, err)
	}
	for _, value := range []string{"bad", " 2026-11-29T00:00:00Z"} {
		if _, err := utils.NormalizeStrictRFC3339UTC(value); !errors.Is(err, utils.ErrInvalidTime) {
			t.Fatalf("accepted invalid time: %q", value)
		}
	}
}
