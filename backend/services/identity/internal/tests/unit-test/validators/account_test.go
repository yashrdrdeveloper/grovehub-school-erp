package validator_test

import (
	"testing"

	coresearch "grovehub/libs/core/search"
	validator "grovehub/services/identity/internal/app/validators/account/impl"
)

func TestAccountUsernamePolicy(t *testing.T) {
	rules := validator.New()
	username, err := rules.NormalizeUsername(" John.Doe ")
	if err != nil || username != "john.doe" {
		t.Fatal("username normalization failed", err)
	}
	if _, err := rules.NormalizeUsername("bad username"); err == nil {
		t.Fatal("space should be rejected")
	}
}

func TestAccountSearchRejectsUntrustedFields(t *testing.T) {
	rules := validator.New()
	_, err := rules.Search(coresearch.Request{
		Filters: []coresearch.Filter{{
			Field: "passwordHash",
			Operator: coresearch.OperatorEquals,
			Values: []string{"bad"},
		}},
	})
	if err == nil {
		t.Fatal("passwordHash filter accepted")
	}
	_, err = rules.Search(coresearch.Request{
		Filters: []coresearch.Filter{{
			Field: "status",
			Operator: coresearch.OperatorEquals,
			Values: []string{"ACTIVE OR 1=1"},
		}},
	})
	if err == nil {
		t.Fatal("untrusted status accepted")
	}
}
