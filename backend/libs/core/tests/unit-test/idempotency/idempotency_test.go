package idempotency_test

import (
	"testing"

	coreidempotency "grovehub/libs/core/idempotency"
)

func TestValidateKey(t *testing.T) {
	key, err := coreidempotency.ValidateKey("  register:019-test  ")
	if err != nil {
		t.Fatalf("ValidateKey() error = %v", err)
	}
	if key != "register:019-test" {
		t.Fatalf("ValidateKey() = %q", key)
	}
	if _, err := coreidempotency.ValidateKey("key with spaces"); err != coreidempotency.ErrInvalidKey {
		t.Fatalf("ValidateKey() unsafe key error = %v", err)
	}
}

func TestFingerprintStable(t *testing.T) {
	request := struct {
		Email string `json:"email"`
	}{Email: "user@example.com"}
	left, err := coreidempotency.Fingerprint(request)
	if err != nil {
		t.Fatalf("Fingerprint() error = %v", err)
	}
	right, err := coreidempotency.Fingerprint(request)
	if err != nil {
		t.Fatalf("Fingerprint() second error = %v", err)
	}
	if left != right || len(left) != 64 {
		t.Fatalf("Fingerprint() values = %q, %q", left, right)
	}
}
