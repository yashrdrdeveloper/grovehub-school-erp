package operatoraudit

import (
	"context"
	"errors"
	"net/http"

	core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)

const OperatorAuditRoute = "POST /audit/restricted/operator-activities"

const OperatorAuditScope = "internal.audit.operator.write"

func New(keys, other []core.WorkloadKeyConfig) (core.WorkloadVerifier, error) {
	for _, key := range keys {
		if key.Subject != "internal-tool-audit" ||
			key.Issuer != "internal-tool-audit" ||
			key.AllowDelegation ||
			len(key.Scopes) != 1 ||
			key.Scopes[0] != OperatorAuditScope {
			return nil, errors.New("invalid operator audit workload key")
		}
		for _, k := range other {
			if k.KeyID == key.KeyID || k.PublicKey == key.PublicKey {
				return nil, errors.New("operator audit key must be independent")
			}
		}
	}

	credentials, err := core.NewEd25519Credentials(keys)
	if err != nil {
		return nil, err
	}

	return &verifier{scoped: core.ScopedWorkloadVerifier{Credentials: credentials, Audience: "platform", Routes: map[string]core.RoutePolicy{OperatorAuditRoute: {Scope: OperatorAuditScope}}}}, nil
}

func (v *verifier) VerifyWorkload(ctx context.Context, r *http.Request) (context.Context, error) {
	if r.Method != "POST" ||
		r.URL == nil ||
		r.URL.Path != "/audit/restricted/operator-activities" ||
		r.URL.RawQuery != "" ||
		r.URL.ForceQuery ||
		r.URL.RawPath != "" ||
		len(r.Header.Values("Cookie")) != 0 ||
		len(r.Header.Values("X-GH-Identity-Proof")) != 0 {
		return nil, core.ErrWorkloadDenied
	}

	verified, err := v.scoped.VerifyWorkload(ctx, r)
	if err != nil {
		return nil, err
	}

	claims, ok := core.WorkloadFrom(verified)
	if !ok ||
		claims.Subject != "internal-tool-audit" ||
		len(claims.Scopes) != 1 ||
		claims.Scopes[0] != OperatorAuditScope ||
		claims.ActorGID != "" ||
		claims.TenantID != "" ||
		claims.SessionID != "" {
		return nil, core.ErrWorkloadDenied
	}

	return verified, nil
}

// @Strategy
type verifier struct{ scoped core.ScopedWorkloadVerifier }

var _ core.WorkloadVerifier = (*verifier)(nil)
