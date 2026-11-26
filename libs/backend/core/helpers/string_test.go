package helpers

import "testing"

func TestOptionalString(t *testing.T) {
    if OptionalString("") != nil { t.Fatal("empty input must map to nil") }
    got := OptionalString("  value  ")
    if got == nil || *got != "  value  " { t.Fatalf("unexpected optional string: %v", got) }
}

func TestNullableString(t *testing.T) {
    if got := NullableString(nil); got != "" { t.Fatalf("nil should map to empty, got %q", got) }
    value := "  value  "
    if got := NullableString(&value); got != value { t.Fatalf("expected unchanged value, got %q", got) }
}
