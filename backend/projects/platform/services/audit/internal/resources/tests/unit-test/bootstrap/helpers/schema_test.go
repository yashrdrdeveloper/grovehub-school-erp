package helpers_test

import (
	"errors"
	"testing"

	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	"grovehub/services/audit/internal/app/bootstrap/constants"
	"grovehub/services/audit/internal/app/bootstrap/helpers"
)

type schemaRepository struct {
	failAt int
	calls  int
	cause  error
}

func (r *schemaRepository) check() error {
	r.calls++
	if r.calls == r.failAt {
		return r.cause
	}
	return nil
}

func (r *schemaRepository) CheckActivity() error      { return r.check() }
func (r *schemaRepository) CheckRejection() error     { return r.check() }
func (r *schemaRepository) CheckReceipt() error       { return r.check() }
func (r *schemaRepository) CheckOperator() error      { return r.check() }
func (r *schemaRepository) CheckSchoolReceipt() error { return r.check() }

func TestCheckSchema(t *testing.T) {
	cause := errors.New("missing column")
	definitions := []coreerrors.Definition{
		constants.ActivityMigrationRequired,
		constants.RejectionMigrationRequired,
		constants.ReceiptMigrationRequired,
		constants.OperatorMigrationRequired,
		constants.SchoolReceiptMigrationRequired,
	}

	for i, definition := range definitions {
		t.Run(string(definition.Code()), func(t *testing.T) {
			repository := &schemaRepository{failAt: i + 1, cause: cause}
			err := helpers.CheckSchema(repository)
			var application *coreerrors.Error
			if !errors.As(err, &application) || application.Definition().Code() != definition.Code() {
				t.Fatalf("wrong error: %v", err)
			}
			if !errors.Is(err, cause) || repository.calls != i+1 {
				t.Fatalf("schema check must stop at the first failure: %v, %d", err, repository.calls)
			}
		})
	}

	repository := &schemaRepository{}
	if err := helpers.CheckSchema(repository); err != nil || repository.calls != len(definitions) {
		t.Fatalf("schema check: %v, calls=%d", err, repository.calls)
	}
}
