package rest

import (
	"encoding/json"
	"net/http"
)

func NoContent() Result[Empty] {
	return Result[Empty]{StatusCode: http.StatusNoContent}
}

// RawWriter is for endpoints with unwrapped JSON success and status-only errors.
func RawWriter[T any]() Writer[T] {
	return Writer[T]{
		Success: func(w http.ResponseWriter, _ *http.Request, result Result[T]) error {
			for name, values := range result.Headers {
				for _, value := range values {
					w.Header().Add(name, value)
				}
			}
			status := result.StatusCode
			if status == 0 {
				status = http.StatusOK
			}
			if status == http.StatusNoContent {
				w.WriteHeader(status)
				return nil
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			return json.NewEncoder(w).Encode(result.Data)
		},
		Error: func(w http.ResponseWriter, _ *http.Request, failure Failure) error {
			w.WriteHeader(failure.StatusCode)
			return nil
		},
	}
}
