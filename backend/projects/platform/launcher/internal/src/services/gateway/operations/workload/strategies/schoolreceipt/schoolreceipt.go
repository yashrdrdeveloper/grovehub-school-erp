package schoolreceipt

import (
	"context"
	"errors"
	core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/helpers"
	"net/http"
	"strings"
)

const ReceiptRoute = "POST /audit/restricted/school-rename-receipts"

const ReceiptScope = "internal.audit.school.rename-receipt.write"

const ReceiptIssuer = "school-registry-audit"

func New(keys, other []core.WorkloadKeyConfig) (core.WorkloadVerifier, error) {
	for i, k := range keys {
		if k.Issuer != ReceiptIssuer || !strings.HasPrefix(k.Subject, ReceiptIssuer+":") || !helpers.ValidSchoolID(strings.TrimPrefix(k.Subject, ReceiptIssuer+":")) || k.AllowDelegation || len(k.Scopes) != 1 || k.Scopes[0] != ReceiptScope {
			return nil, errors.New("invalid school receipt key")
		}
		for _, old := range append(append([]core.WorkloadKeyConfig(nil), other...), keys[:i]...) {
			if old.KeyID == k.KeyID || old.PublicKey == k.PublicKey {
				return nil, errors.New("school receipt key must be independent")
			}
		}
	}

	credentials, err := core.NewEd25519Credentials(keys)
	if err != nil {
		return nil, err
	}

	return &verifier{scoped: core.ScopedWorkloadVerifier{Credentials: credentials, Audience: "platform", Routes: map[string]core.RoutePolicy{ReceiptRoute: {Scope: ReceiptScope}}}}, nil
}

func (v *verifier) VerifyWorkload(ctx context.Context, r *http.Request) (context.Context, error) {
	if r.Method != "POST" || r.URL == nil || r.URL.Path != "/audit/restricted/school-rename-receipts" || r.URL.RawPath != "" || r.URL.RawQuery != "" || r.URL.ForceQuery || len(r.Header.Values("Cookie")) != 0 || len(r.Header.Values("X-GH-Identity-Proof")) != 0 {
		return nil, core.ErrWorkloadDenied
	}

	verified, err := v.scoped.VerifyWorkload(ctx, r)
	if err != nil {
		return nil, err
	}

	claims, ok := core.WorkloadFrom(verified)
	if !ok || !strings.HasPrefix(claims.Subject, ReceiptIssuer+":") || !helpers.ValidSchoolID(strings.TrimPrefix(claims.Subject, ReceiptIssuer+":")) || len(claims.Scopes) != 1 || claims.Scopes[0] != ReceiptScope || claims.ActorGID != "" || claims.TenantID != "" || claims.SessionID != "" {
		return nil, core.ErrWorkloadDenied
	}

	return verified, nil
}

// @Strategy
type verifier struct{ scoped core.ScopedWorkloadVerifier }

var _ core.WorkloadVerifier = (*verifier)(nil)
