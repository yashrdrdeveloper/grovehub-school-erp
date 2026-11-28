package workload

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/internalcatalog"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/operatoraudit"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/operatorauditread"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/schoolexecution"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/schoolreceipt"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/schoolreceiptread"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/schoolstatusreceipt"
	"grovehub/launcher/internal/src/services/gateway/policies"
	school "grovehub/libs/platform/session"
	mail "grovehub/services/integration/client/mail"
	integration "grovehub/services/integration/client/principal"
)

const gmailWebhook = "POST /integrations/restricted/webhooks/gmail"

type delegatedIdentity interface {
	VerifyDelegatedSession(context.Context, string) (school.Principal, error)
	school.ProtectedRouteAuthorizer
	school.AccessVerifier
}

// @Service
type platformWorkloads struct {
	scoped              core.ScopedWorkloadVerifier
	internalCatalog     core.WorkloadVerifier
	operatorAudit       core.WorkloadVerifier
	operatorAuditRead   core.WorkloadVerifier
	schoolExecution     core.WorkloadVerifier
	schoolReceipt       core.WorkloadVerifier
	schoolStatusReceipt core.WorkloadVerifier
	receiptRead         core.WorkloadVerifier
	policies            map[string]core.RoutePolicy
	identity            delegatedIdentity
	webhook             core.WorkloadVerifier
	tenantCheck         school.TenantCheck
}

type WorkloadConfig struct {
	Keys                    []core.WorkloadKeyConfig
	InternalToolKeys        []core.WorkloadKeyConfig
	OperatorAuditKeys       []core.WorkloadKeyConfig
	OperatorAuditReadKeys   []core.WorkloadKeyConfig
	SchoolExecutionKeys     []core.WorkloadKeyConfig
	SchoolReceiptKeys       []core.WorkloadKeyConfig
	SchoolStatusReceiptKeys []core.WorkloadKeyConfig
	SchoolReceiptReadKeys   []core.WorkloadKeyConfig
	Routes                  map[string]core.RoutePolicy
}

// NewPlatformWorkloads uses explicit route grants. Empty config grants no
// service-to-service access. The existing Gmail webhook verifier stays isolated.
func NewPlatformWorkloads(raw string, identity delegatedIdentity, webhook core.WorkloadVerifier, checks ...school.TenantCheck) (core.WorkloadVerifier, error) {
	v := &platformWorkloads{identity: identity, webhook: webhook}
	if len(checks) > 0 {
		v.tenantCheck = checks[0]
	}

	if strings.TrimSpace(raw) == "" {
		return v, nil
	}

	var cfg WorkloadConfig
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("invalid workload configuration")
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, errors.New("invalid trailing workload configuration")
	}

	if len(cfg.InternalToolKeys) > 0 {
		catalog, err := internalcatalog.New(cfg.InternalToolKeys, cfg.Keys)
		if err != nil {
			return nil, err
		}

		v.internalCatalog = catalog
	}

	if len(cfg.OperatorAuditKeys) > 0 {
		other := append(append([]core.WorkloadKeyConfig(nil), cfg.Keys...), cfg.InternalToolKeys...)
		audit, err := operatoraudit.New(cfg.OperatorAuditKeys, other)
		if err != nil {
			return nil, err
		}

		v.operatorAudit = audit
	}

	if len(cfg.OperatorAuditReadKeys) > 0 {
		other := append(append(append([]core.WorkloadKeyConfig(nil), cfg.Keys...), cfg.InternalToolKeys...), cfg.OperatorAuditKeys...)
		reader, err := operatorauditread.New(cfg.OperatorAuditReadKeys, other)
		if err != nil {
			return nil, err
		}

		v.operatorAuditRead = reader
	}

	if len(cfg.SchoolExecutionKeys) > 0 {
		other := append(append(append(append([]core.WorkloadKeyConfig(nil), cfg.Keys...), cfg.InternalToolKeys...), cfg.OperatorAuditKeys...), cfg.OperatorAuditReadKeys...)
		school, err := schoolexecution.New(cfg.SchoolExecutionKeys, other)
		if err != nil {
			return nil, err
		}

		v.schoolExecution = school
	}

	if len(cfg.SchoolReceiptKeys) > 0 {
		other := append(append(append(append(append([]core.WorkloadKeyConfig(nil), cfg.Keys...), cfg.InternalToolKeys...), cfg.OperatorAuditKeys...), cfg.OperatorAuditReadKeys...), cfg.SchoolExecutionKeys...)
		receipt, err := schoolreceipt.New(cfg.SchoolReceiptKeys, other)
		if err != nil {
			return nil, err
		}

		v.schoolReceipt = receipt
	}

	if len(cfg.SchoolReceiptReadKeys) > 0 {
		other := append(append(append(append(append(append([]core.WorkloadKeyConfig(nil), cfg.Keys...), cfg.InternalToolKeys...), cfg.OperatorAuditKeys...), cfg.OperatorAuditReadKeys...), cfg.SchoolExecutionKeys...), cfg.SchoolReceiptKeys...)
		reader, err := schoolreceiptread.New(cfg.SchoolReceiptReadKeys, other)
		if err != nil {
			return nil, err
		}

		v.receiptRead = reader
	}

	if len(cfg.SchoolStatusReceiptKeys) > 0 {
		other := append(append(append(append(append(append(append([]core.WorkloadKeyConfig(nil), cfg.Keys...), cfg.InternalToolKeys...), cfg.OperatorAuditKeys...), cfg.OperatorAuditReadKeys...), cfg.SchoolExecutionKeys...), cfg.SchoolReceiptKeys...), cfg.SchoolReceiptReadKeys...)
		status, err := schoolstatusreceipt.New(cfg.SchoolStatusReceiptKeys, other)
		if err != nil {
			return nil, err
		}

		v.schoolStatusReceipt = status
	}

	if len(cfg.Keys) == 0 &&
		len(cfg.Routes) == 0 &&
		(len(cfg.InternalToolKeys) > 0 || len(cfg.OperatorAuditKeys) > 0 || len(cfg.OperatorAuditReadKeys) > 0 || len(cfg.SchoolExecutionKeys) > 0 || len(cfg.SchoolReceiptKeys) > 0 || len(cfg.SchoolReceiptReadKeys) > 0 || len(cfg.SchoolStatusReceiptKeys) > 0) {
		return v, nil
	}
	// Mail diagnostics require a dedicated single-scope, nondelegating operator key.
	for _, key := range cfg.Keys {
		for _, scope := range key.Scopes {
			if scope == mail.ReconciliationScope ||
				key.Subject == mail.ReconciliationSubject ||
				key.Issuer == mail.ReconciliationSubject {
				if key.Subject != mail.ReconciliationSubject ||
					key.Issuer != mail.ReconciliationSubject ||
					key.AllowDelegation ||
					len(key.Scopes) != 1 ||
					scope != mail.ReconciliationScope {
					return nil, errors.New("invalid mail reconciliation workload key")
				}
			}
		}
	}

	credentials, err := core.NewEd25519Credentials(cfg.Keys)
	if err != nil {
		return nil, err
	}

	if len(cfg.Routes) == 0 {
		return nil, errors.New("no workload route policies")
	}

	for pattern, policy := range cfg.Routes {
		if !policies.Validate(pattern, policy) {
			return nil, fmt.Errorf("invalid workload route policy: %s", pattern)
		}
		if policy.RequireActor && identity == nil {
			return nil, errors.New("delegated route has no supported identity authorization boundary")
		}
	}
	// The proof-based School boundary independently enforces actor presence.
	// The generic verifier must not require an unverified actor assertion.
	routes := make(map[string]core.RoutePolicy, len(cfg.Routes))
	for pattern, policy := range cfg.Routes {
		routes[pattern] = core.RoutePolicy{Scope: policy.Scope}
	}

	v.scoped = core.ScopedWorkloadVerifier{Credentials: credentials, Audience: "platform", Routes: routes}
	v.policies = cfg.Routes
	return v, nil
}

func (v *platformWorkloads) VerifyWorkload(ctx context.Context, r *http.Request) (context.Context, error) {
	if r == nil {
		return nil, core.ErrWorkloadDenied
	}

	if r.Pattern == gmailWebhook {
		if v.webhook == nil {
			return nil, core.ErrWorkloadDenied
		}
		return v.webhook.VerifyWorkload(ctx, r)
	}

	if r.Pattern == schoolstatusreceipt.StatusReceiptRoute {
		return verifyStrategy(v.schoolStatusReceipt, ctx, r)
	}

	if r.Pattern == schoolreceiptread.ReceiptReadRoute {
		return verifyStrategy(v.receiptRead, ctx, r)
	}

	if r.Pattern == schoolreceipt.ReceiptRoute {
		return verifyStrategy(v.schoolReceipt, ctx, r)
	}

	if r.Pattern == schoolexecution.SchoolAuditCheckRoute {
		return verifyStrategy(v.schoolExecution, ctx, r)
	}

	if r.Pattern == operatorauditread.OperatorAuditReadRoute {
		return verifyStrategy(v.operatorAuditRead, ctx, r)
	}

	if r.Pattern == operatoraudit.OperatorAuditRoute {
		return verifyStrategy(v.operatorAudit, ctx, r)
	}

	if r.Pattern == internalcatalog.InternalCatalogRoute {
		if verified, err := verifyStrategy(v.internalCatalog, ctx, r); err == nil {
			return verified, nil
		}
	}

	ctx = school.WithTenantCheck(ctx, v.tenantCheck)
	verified, err := v.scoped.VerifyWorkload(ctx, r)
	if err != nil {
		return nil, err
	}

	claims, authenticated := core.WorkloadFrom(verified)
	policy, granted := v.policies[r.Pattern]
	if !authenticated || !granted || len(claims.Scopes) != 1 || claims.Scopes[0] != policy.Scope {
		return nil, core.ErrWorkloadDenied
	}

	if r.Pattern == "POST "+mail.ReconciliationPath {
		if policy.RequireActor ||
			claims.Subject != mail.ReconciliationSubject ||
			claims.ActorGID != "" ||
			claims.TenantID != "" ||
			claims.SessionID != "" ||
			r.URL == nil ||
			r.URL.Path != mail.ReconciliationPath ||
			r.URL.RawPath != "" ||
			r.URL.RawQuery != "" ||
			r.URL.ForceQuery ||
			len(r.Header.Values("Cookie")) != 0 ||
			len(r.Header.Values("X-GH-Identity-Proof")) != 0 {
			return nil, core.ErrWorkloadDenied
		}
		return verified, nil
	}

	if !policy.RequireActor {
		if (foundationRoute(r.Pattern) && claims.Subject != "school-foundation") ||
			(!foundationRoute(r.Pattern) && claims.Subject != "school") ||
			claims.ActorGID != "" ||
			claims.TenantID != "" ||
			claims.SessionID != "" {
			return nil, core.ErrWorkloadDenied
		}
		return verified, nil
	}

	if v.identity == nil || claims.Subject != "school" {
		return nil, core.ErrWorkloadDenied
	}
	// School must prove possession of the actual access secret. A session ID
	// or a signed assertion alone cannot authenticate a browser user.
	values := r.Header.Values("X-GH-Identity-Proof")
	if len(values) != 1 || len(values[0]) == 0 || len(values[0]) > 2048 || strings.ContainsAny(values[0], " \t\r\n,;") {
		return nil, core.ErrWorkloadDenied
	}

	if claims.SessionID != "" || claims.ActorGID != "" || claims.TenantID != "" {
		return nil, core.ErrWorkloadDenied
	}

	actor, err := v.identity.VerifyAccess(verified, values[0])
	if err != nil || actor.SessionID == "" || actor.AccountID == "" {
		return nil, core.ErrWorkloadDenied
	}
	// Logout and revocation remain available to disabled users.
	switch r.Pattern {
	case "POST /identity/restricted/auth/logout", "POST /identity/restricted/auth/logout-all", "DELETE /identity/restricted/auth/sessions/{sessionId}":
	default:
		if school.CheckTenantSession(verified, actor.TenantID, actor.SessionGeneration) != nil {
			return nil, core.ErrWorkloadDenied
		}
	}

	r.Header.Del("X-GH-Identity-Proof")
	verified = school.WithPrincipal(verified, actor)
	// Reuse existing route policy during migration; unknown routes deny.
	oldPattern := strings.Replace(r.Pattern, "/restricted/", "/protected/", 1)
	if integrationRoute(r.Pattern) {
		// Integration controllers enforce tenant + owner on each resource and
		// independently require admin for system-mail configuration.
		if actor.TenantID == "" || actor.PersonGID == "" {
			return nil, core.ErrWorkloadDenied
		}
	} else if err := v.identity.AuthorizeProtectedRoute(verified, oldPattern); err != nil {
		return nil, core.ErrWorkloadDenied
	}

	verified = integration.WithVerifiedPrincipal(verified, integration.VerifiedPrincipal{TenantID: actor.TenantID, ActorGID: actor.PersonGID})
	return verified, nil
}

func integrationRoute(pattern string) bool {
	switch pattern {
	case "POST /integrations/restricted/connections":
		return true
	case "GET /integrations/restricted/connections/{id}":
		return true
	case "GET /integrations/restricted/connections":
		return true
	case "DELETE /integrations/restricted/connections/{id}":
		return true
	case "POST /integrations/restricted/connections/{id}/oauth/start":
		return true
	case "POST /integrations/restricted/connections/{id}/test":
		return true
	case "GET /integrations/restricted/connections/{id}/tests":
		return true
	case "POST /integrations/restricted/connections/{id}/watch":
		return true
	case "POST /integrations/restricted/connections/system-mail/sender":
		return true
	case "POST /integrations/restricted/connections/{id}/actions/{action}":
		return true
	case "GET /integrations/restricted/providers":
		return true
	case "GET /integrations/restricted/providers/{provider}":
		return true
	default:
		return false
	}
}

func foundationRoute(pattern string) bool {
	switch pattern {
	case "POST /foundation/restricted/persons", "GET /foundation/restricted/persons/{gid}", "PATCH /foundation/restricted/persons/{gid}", "POST /foundation/restricted/persons/search", "PATCH /foundation/restricted/persons/{gid}/deactivate":
		return true
	default:
		return false
	}
}

func verifyStrategy(verifier core.WorkloadVerifier, ctx context.Context, r *http.Request) (context.Context, error) {
	if verifier == nil {
		return nil, core.ErrWorkloadDenied
	}
	return verifier.VerifyWorkload(ctx, r)
}
