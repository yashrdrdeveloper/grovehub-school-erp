package internalcatalog

import (
	"context"
	"errors"
	"net/http"

	core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)

const InternalCatalogRoute = "GET /integrations/restricted/providers"

const InternalCatalogScope = "internal.integration.catalog.read"

func New(keys, schoolKeys []core.WorkloadKeyConfig) (core.WorkloadVerifier, error) {
	for _, key := range keys {
		if key.Subject != "internal-tool" || key.Issuer != "internal-tool" || key.AllowDelegation || len(key.Scopes) != 1 || key.Scopes[0] != InternalCatalogScope {
			return nil, errors.New("invalid internal tool catalog key")
		}
		for _, school := range schoolKeys {
			if school.KeyID == key.KeyID || school.PublicKey == key.PublicKey {
				return nil, errors.New("internal tool key must be independent")
			}
		}
	}

	credentials, err := core.NewEd25519Credentials(keys)
	if err != nil {
		return nil, err
	}

	return &verifier{scoped: core.ScopedWorkloadVerifier{Credentials: credentials, Audience: "platform", Routes: map[string]core.RoutePolicy{InternalCatalogRoute: {Scope: InternalCatalogScope}}}}, nil
}

func (v *verifier) VerifyWorkload(ctx context.Context, r *http.Request) (context.Context, error) {
	if r.Pattern != InternalCatalogRoute || r.Method != "GET" || r.URL == nil || r.URL.Path != "/integrations/restricted/providers" || r.URL.RawQuery != "" || r.URL.RawPath != "" || len(r.Header.Values("Cookie")) != 0 || len(r.Header.Values("X-GH-Identity-Proof")) != 0 {
		return nil, core.ErrWorkloadDenied
	}

	verified, err := v.scoped.VerifyWorkload(ctx, r)
	if err != nil {
		return nil, err
	}

	claims, ok := core.WorkloadFrom(verified)
	if !ok || claims.Subject != "internal-tool" || len(claims.Scopes) != 1 || claims.Scopes[0] != InternalCatalogScope || claims.ActorGID != "" || claims.TenantID != "" || claims.SessionID != "" {
		return nil, core.ErrWorkloadDenied
	}

	return verified, nil
}

// @Strategy
type verifier struct{ scoped core.ScopedWorkloadVerifier }

var _ core.WorkloadVerifier = (*verifier)(nil)
