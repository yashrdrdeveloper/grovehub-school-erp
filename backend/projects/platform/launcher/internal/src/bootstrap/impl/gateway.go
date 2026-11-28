package impl

import (
	"errors"
	"fmt"
	"net/http"

	"grovehub/launcher/internal/src/services/gateway"
	gatewayimpl "grovehub/launcher/internal/src/services/gateway/impl"
	tenantaccess "grovehub/launcher/internal/src/services/gateway/operations/tenant/access"
	"grovehub/launcher/internal/src/services/gateway/operations/workload"
	"grovehub/libs/platform/session"
)

func (b *builder) newTenantState() (session.TenantStateLookup, session.TenantCheck, error) {
	lookup, err := tenantaccess.NewTenantState(tenantaccess.TenantAccessConfig{Origin: b.config.TenantOrigin, KeyID: b.config.TenantKeyID, PrivateKey: b.config.TenantPrivateKey})
	if err != nil {
		return nil, nil, err
	}

	check := session.ActiveTenantCheck(lookup)
	if b.config.TenantRequired && check == nil {
		return nil, nil, errors.New("tenant access verifier required")
	}

	return lookup, check, nil
}

func (b *builder) newGateway(mux *http.ServeMux, modules serviceModules, tenantCheck session.TenantCheck) (gateway.Gateway, error) {
	workloads, err := workload.NewPlatformWorkloads(b.config.Workload, modules.identity, modules.integration, tenantCheck)
	if err != nil {
		return nil, fmt.Errorf("initialize workload gateway: %w", err)
	}

	return gatewayimpl.New(mux, modules.identity, workloads, b.config.AllowedOrigin, nil, modules.identity.RouteLimits()), nil
}
