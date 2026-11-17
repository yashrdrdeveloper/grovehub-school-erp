package credential_test

import (
	"testing"

	"grovehub/services/identity/internal/app/services/credential/helpers"
)

func TestArgon2HashAndVerify(t *testing.T) {
	password := "correct horse battery staple"
	hash, err := helpers.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	if hash == password {
		t.Fatal("plaintext credential was retained")
	}

	ok, err := helpers.VerifyPassword(password, hash)
	if err != nil || !ok {
		t.Fatalf("valid credential rejected: %v", err)
	}
	ok, err = helpers.VerifyPassword("wrong-password", hash)
	if err != nil || ok {
		t.Fatal("invalid password was accepted")
	}
}
