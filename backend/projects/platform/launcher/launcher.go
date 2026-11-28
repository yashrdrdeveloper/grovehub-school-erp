package launcher

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"gorm.io/gorm"

	coregateway "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	"grovehub/launcher/internal/src/bootstrap"
	bootstrapimpl "grovehub/launcher/internal/src/bootstrap/impl"
	commandbootstrap "grovehub/launcher/internal/src/operations/command/bootstrap"
	"grovehub/launcher/internal/src/operations/command/registry"
	external "grovehub/launcher/internal/src/services/gateway/operations/external/impl"
	tenantaccess "grovehub/launcher/internal/src/services/gateway/operations/tenant/access"
	"grovehub/launcher/internal/src/services/gateway/operations/workload"
	"grovehub/libs/platform/session"
)

func Run(logger *slog.Logger) error {
	return bootstrap.Run(bootstrapimpl.New(logger))
}

// Operator commands share the same Platform-owned database and logger setup.
func NewLogger(output io.Writer) (*slog.Logger, error) {
	return bootstrapimpl.NewLogger(output)
}

func NewPostgresDB(ctx context.Context) (*gorm.DB, error) {
	return bootstrapimpl.NewPostgresDB(ctx)
}

func NewAuditPostgresDB(ctx context.Context, shared *gorm.DB) (*gorm.DB, func(), error) {
	return bootstrapimpl.NewAuditPostgresDB(ctx, shared)
}

var ErrDatabaseURLRequired = bootstrapimpl.ErrDatabaseURLRequired

// Public composition contracts alias the Platform-private session boundary.
type Principal = session.Principal

type AccessVerifier = session.AccessVerifier

type ProtectedRouteAuthorizer = session.ProtectedRouteAuthorizer

type TenantCheck = session.TenantCheck

type TenantState = session.TenantState

type TenantStateLookup = session.TenantStateLookup

// DelegatedIdentity is the verified actor boundary for composing the gateway.
type DelegatedIdentity interface {
	VerifyDelegatedSession(context.Context, string) (session.Principal, error)
	session.ProtectedRouteAuthorizer
	session.AccessVerifier
}

type TenantAccessConfig = tenantaccess.TenantAccessConfig

func NewTenantState(c TenantAccessConfig) (session.TenantStateLookup, error) {
	return tenantaccess.NewTenantState(c)
}

func NewWorkloadVerifier(raw string, identity DelegatedIdentity, webhook coregateway.WorkloadVerifier, checks ...session.TenantCheck) (coregateway.WorkloadVerifier, error) {
	return workload.NewPlatformWorkloads(raw, identity, webhook, checks...)
}

func NewExternalGateway(mux *http.ServeMux, identity session.AccessVerifier, workloads coregateway.WorkloadVerifier, origin string, limits ...map[string]int) http.Handler {
	return external.New(mux, identity, workloads, origin, limits...)
}

// Execute dispatches the requested operation without starting unrelated services.
func Execute(args []string, output, diagnostics io.Writer) int {
	return commandbootstrap.Dispatch(args, output, diagnostics, func() (registry.Registry, error) {
		return commandbootstrap.New(Run)
	})
}
