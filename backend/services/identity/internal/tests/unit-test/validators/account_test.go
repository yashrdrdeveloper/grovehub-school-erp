package validator_test

import (
	"testing"
	"reflect"
	"strings"

	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	searchconfig "grovehub/services/identity/internal/app/repository/account/helpers/search"

	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
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

func TestAccountSearchAcceptsEntityStatuses(t *testing.T) {
	rules := validator.New()
	for _, status := range []entity.AccountStatus{entity.Active, entity.PendingSetup, entity.Suspended, entity.Deactivated} {
		_, err := rules.Search(coresearch.Request{
			Filters: []coresearch.Filter{{
				Field: accountdto.FieldStatus,
				Operator: coresearch.OperatorEquals,
				Values: []string{string(status)},
			}},
		})
		if err != nil {
			t.Errorf("entity status %s rejected: %v", status, err)
		}
	}
}

func TestAccountSearchFieldsMatchDTOAndSQL(t *testing.T) {
	fields := map[string]string{
		accountdto.FieldAccountID: "AccountID",
		accountdto.FieldPersonGID: "PersonGID",
		accountdto.FieldUsername: "Username",
		accountdto.FieldEmail: "Email",
		accountdto.FieldStatus: "Status",
		accountdto.FieldCreatedAt: "CreatedAt",
		accountdto.FieldUpdatedAt: "UpdatedAt",
	}
	typeOfAccount := reflect.TypeOf(accountdto.Account{})
	for fieldName, goName := range fields {
		field, ok := typeOfAccount.FieldByName(goName)
		if !ok || strings.Split(field.Tag.Get("json"), ",")[0] != fieldName {
			t.Errorf("search field %q does not match response JSON for %s", fieldName, goName)
		}
		if _, ok := searchconfig.AccountConfig.ScalarFields[fieldName]; !ok {
			t.Errorf("search field %q missing from repository SQL allowlist", fieldName)
		}
	}
	for _, fieldName := range []string{accountdto.FieldAccountID, accountdto.FieldUsername, accountdto.FieldStatus, accountdto.FieldCreatedAt, accountdto.FieldUpdatedAt} {
		if _, ok := searchconfig.AccountConfig.SortFields[fieldName]; !ok {
			t.Errorf("sort field %q missing from repository SQL allowlist", fieldName)
		}
	}
}
