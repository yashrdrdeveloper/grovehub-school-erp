package impl

import (
	"context"
	"log/slog"
	"net/http"
	"sync"

	schoolgateway "grovehub/libs/schoolgateway"
	accountclient "grovehub/services/identity/client/account"
	"grovehub/services/identity/internal/app/bootstrap"
	security "grovehub/services/identity/internal/app/services/security"
)

type runtime struct {
	logger   *slog.Logger
	routes   controllerGraph
	client   accountclient.Client
	verifier security.AccessVerifier
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
