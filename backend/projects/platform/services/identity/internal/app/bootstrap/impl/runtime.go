package impl

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	schoolgateway "grovehub/launcher/services/gateway/session"
	accountclient "grovehub/services/identity/client/account"
	"grovehub/services/identity/internal/app/bootstrap"
	security "grovehub/services/identity/internal/app/services/security"
	"grovehub/services/identity/internal/app/services/security/operations/routepolicy"
)

type runtime struct {
	logger   *slog.Logger
	routes   controllerGraph
	client   accountclient.Client
	verifier security.AccessVerifier
	policy *routepolicy.Policy
	stop     context.CancelFunc
	done     chan struct{}
	once     sync.Once
}

var _ bootstrap.Runtime = (*runtime)(nil)

func (module *runtime) RegisterHTTP(mux *http.ServeMux) {
	module.routes.RegisterHTTP(mux)
}

func (module *runtime) VerifyAccess(ctx context.Context, secret string) (schoolgateway.Principal, error) {
	return module.verifier.VerifyAccess(ctx, secret)
}

func (module *runtime) AccountClient() accountclient.Client {
	return module.client
}

func (module *runtime) Close() {
	module.once.Do(func() {
		if module.stop != nil {
			module.stop()
		}
		if module.done != nil {
			<-module.done
		}
	})
}

func (module *runtime) AuthorizeProtectedRoute(ctx context.Context, pattern string) error {
	return module.policy.Authorize(ctx, pattern)
}

func (module *runtime) VerifyDelegatedSession(ctx context.Context, sessionID string) (schoolgateway.Principal, error) {
 verifier, ok := module.verifier.(interface {
  VerifyDelegatedSession(context.Context,string)(schoolgateway.Principal,error)
 })
 if !ok { return schoolgateway.Principal{}, errors.New("delegated session unavailable") }
 return verifier.VerifyDelegatedSession(ctx, sessionID)
}
