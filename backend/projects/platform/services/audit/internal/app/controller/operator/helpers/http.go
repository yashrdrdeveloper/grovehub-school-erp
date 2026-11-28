package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/core/gateway"
	"grovehub/libs/core/pipeline"
	definitions "grovehub/services/audit/internal/app/constants/error"
	"grovehub/services/audit/internal/app/controller/operator/constants"
	repository "grovehub/services/audit/internal/app/repository/operator"
)

func Authorized(r *http.Request, subject, scope string) bool {
	claims, ok := gateway.WorkloadFrom(r.Context())
	return ok && (subject == "" || claims.Subject == subject) &&
		len(claims.Scopes) == 1 && claims.Scopes[0] == scope &&
		claims.ActorGID == "" && claims.TenantID == "" && claims.SessionID == ""
}

func Decode(w http.ResponseWriter) func(*http.Request, any) error {
	return func(r *http.Request, value any) error {
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, constants.MaxBodyBytes))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(value); err != nil {
			return err
		}

		var extra any
		if decoder.Decode(&extra) != io.EOF {
			return coreerrors.NewError(definitions.OperatorInvalid)
		}
		return nil
	}
}

func Failure(err error) error {
	if err == nil {
		return nil
	}
	var application *coreerrors.Error
	if errors.As(err, &application) {
		return err
	}
	if errors.Is(err, repository.ErrConflict) {
		return coreerrors.NewError(definitions.OperatorConflict)
	}
	return coreerrors.NewError(definitions.OperatorUnavailable)
}

// Writer preserves the existing restricted API: bare statuses and unwrapped JSON.
func Writer[T any]() pipeline.Writer[T] {
	return pipeline.Writer[T]{
		Success: func(w http.ResponseWriter, _ *http.Request, result pipeline.Result[T]) error {
			if result.StatusCode == http.StatusNoContent {
				w.WriteHeader(http.StatusNoContent)
				return nil
			}
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(result.Data)
		},
		Error: func(w http.ResponseWriter, _ *http.Request, failure pipeline.Failure) error {
			w.WriteHeader(failure.StatusCode)
			return nil
		},
	}
}
