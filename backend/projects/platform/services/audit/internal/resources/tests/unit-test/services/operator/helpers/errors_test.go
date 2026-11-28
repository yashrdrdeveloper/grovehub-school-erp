package helpers_test

import (
	"errors"
	"net/http"
	"testing"

	coreerrors "grovehub/libs/core/errors"
	repository "grovehub/services/audit/internal/app/repository/operator"
	"grovehub/services/audit/internal/app/services/operator/helpers"
)

func TestStorageError(t *testing.T) {
	if helpers.StorageError(nil) != nil {
		t.Fatal("nil error changed")
	}
	for _, tc := range []struct {
		cause error
		status int
	}{
		{repository.ErrConflict, http.StatusConflict},
		{errors.New("database connection failure"), http.StatusServiceUnavailable},
	} {
		err := helpers.StorageError(tc.cause)
		var app *coreerrors.Error
		if !errors.As(err, &app) || app.Definition().HTTPStatus() != tc.status || !errors.Is(err, tc.cause) {
			t.Fatalf("wrong storage error: %v", err)
		}
	}
}
