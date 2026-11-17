package authentication_test

import (
	"testing"

	"grovehub/services/identity/internal/app/services/authentication/helpers"
)

func TestTokenSecretsAreRandomAndHashed(t *testing.T) {
	first, err := helpers.NewSecret(32)
	if err != nil {
		t.Fatal(err)
	}
	second, err := helpers.NewSecret(32)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("random tokens collide")
	}

	hash, err := helpers.Digest(first, 32)
	if err != nil || len(hash) != 64 || hash == first {
		t.Fatal("invalid digest")
	}
	if _, err := helpers.Digest(first, 48); err == nil {
		t.Fatal("incorrect token length accepted")
	}
}
