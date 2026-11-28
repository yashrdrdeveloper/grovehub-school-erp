package helpers

import (
	"context"
	"net/http"
	"time"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/core/gateway"
	corehttp "grovehub/libs/core/http"
	"grovehub/libs/core/pipeline"
	corerequest "grovehub/libs/core/request"
	definitions "grovehub/services/audit/internal/app/constants/error"
	"grovehub/services/audit/internal/app/controller/operator/constants"
)

type Policy struct {
	Subject string
	Scope   string
	Timeout time.Duration
}

// Authorized enforces the Audit-specific nondelegated workload policy.
func Authorized(r *http.Request, subject, scope string) bool {
	claims, ok := gateway.WorkloadFrom(r.Context())
	return ok && (subject == "" || claims.Subject == subject) &&
		len(claims.Scopes) == 1 && claims.Scopes[0] == scope &&
		claims.ActorGID == "" && claims.TenantID == "" && claims.SessionID == ""
}

// ValidateReceipt checks the decoded school identifier against the authenticated issuer.
func ValidateReceipt[T any](r *http.Request, validate func(T) error, schoolID func(T) string, issuer, scope string) func(T) error {
	return func(input T) error {
		if err := validate(input); err != nil {
			return err
		}
		if !Authorized(r, issuer+":"+schoolID(input), scope) {
			return coreerrors.NewError(definitions.OperatorForbidden)
		}
		return nil
	}
}

// Run configures the common restricted Audit transport without duplicating the REST pipeline.
func Run[Q corerequest.Contract, R any](w http.ResponseWriter, r *http.Request, policy Policy, config pipeline.StructConfig[Q, R]) {
	w.Header().Set("Cache-Control", "no-store")
	if !Authorized(r, policy.Subject, policy.Scope) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	config.Decode = func(request *http.Request, destination any) error {
		request.Body = http.MaxBytesReader(w, request.Body, constants.MaxBodyBytes)
		return corehttp.DecodeStrictJSON(request, destination)
	}
	config.InvalidRequestBody = definitions.OperatorInvalid
	config.Fallback = definitions.OperatorUnavailable
	config.Writer = pipeline.RawWriter[R]()

	timeout := policy.Timeout
	if timeout == 0 {
		timeout = constants.RequestTimeout
	}
	service := config.Service
	config.Service = func(ctx context.Context, input Q) (pipeline.Result[R], error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return service(ctx, input)
	}
	pipeline.RunStruct(w, r, config)
}
