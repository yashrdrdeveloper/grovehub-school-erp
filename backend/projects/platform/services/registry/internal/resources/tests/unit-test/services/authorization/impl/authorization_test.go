package impl_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"testing"
	"time"

	coregateway "grovehub/libs/core/gateway"
	"grovehub/libs/platform/scope"
	"grovehub/services/registry/client/authorization"
	"grovehub/services/registry/internal/app/services/authorization/impl"
)

const projectID = "10000000-0000-4000-8000-000000000001"
const callerID = "10000000-0000-4000-8000-000000000002"
const tenantID = "20000000-0000-4000-8000-000000000001"
const workloadID = "30000000-0000-4000-8000-000000000001"

type repository struct {
	value authorization.Credential
	state authorization.TenantState
	err error
	calls int
}

func (r *repository) Resolve(_ context.Context, keyID string, target scope.Target, capability string, inactive bool) (authorization.Credential, error) {
	r.calls++
	if keyID != "key-one" || target.Level != scope.Tenant || target.ProjectID != projectID || target.TenantID != tenantID || capability != "audit.read" || inactive {
		return authorization.Credential{}, scope.ErrDenied
	}
	return r.value, r.err
}

func (r *repository) TenantState(_ context.Context, _, _ string) (authorization.TenantState, error) {
	return r.state, r.err
}

func TestLiveVerification(t *testing.T) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil { t.Fatal(err) }
	now := time.Now()
	repo := &repository{value: authorization.Credential{
		WorkloadID: workloadID, CallerProjectID: callerID, Issuer: "issuer",
		Subject: "executor", PublicKey: public,
		ValidFrom: now.Add(-time.Hour), ValidUntil: now.Add(time.Hour),
		WorkloadUpdatedAt: now.Add(-time.Hour),
	}}
	service := impl.NewAuthorizationService(repo, "platform-dev")
	claims := coregateway.WorkloadClaims{
		Subject: "executor", Audience: "platform-dev", ExpiresAt: now.Add(time.Minute),
		Scopes: []string{"audit.read"}, TargetScope: string(scope.Tenant),
		TargetProjectID: projectID, TargetTenantID: tenantID,
		RequestMethod: "POST", RequestURI: "/audit/restricted/activities/search",
	}
	sign := func(value coregateway.WorkloadClaims) string {
		t.Helper()
		token, err := coregateway.SignWorkload("key-one", "issuer", private, value)
		if err != nil { t.Fatal(err) }
		return token
	}
	request := authorization.Request{
		Token: sign(claims), Target: scope.Target{Level: scope.Tenant, ProjectID: projectID, TenantID: tenantID},
		Capability: "audit.read", Method: claims.RequestMethod, RequestURI: claims.RequestURI,
	}
	ctx := context.Background()
	verified, err := service.Verify(ctx, request)
	if err != nil { t.Fatal(err) }
	value, err := scope.Require(verified, request.Target, request.Capability)
	if err != nil || value.CallerProjectID != callerID || value.WorkloadID != workloadID || value.PrincipalType != scope.Workload {
		t.Fatalf("wrong execution identity: %+v %v", value, err)
	}
	for _, change := range []func(*coregateway.WorkloadClaims){
		func(v *coregateway.WorkloadClaims) { v.Audience = "platform-prod" },
		func(v *coregateway.WorkloadClaims) { v.TargetScope = "" },
		func(v *coregateway.WorkloadClaims) { v.TargetProjectID = callerID },
		func(v *coregateway.WorkloadClaims) { v.TargetTenantID = projectID },
		func(v *coregateway.WorkloadClaims) { v.RequestMethod = "DELETE" },
		func(v *coregateway.WorkloadClaims) { v.RequestURI += "?forged=true" },
		func(v *coregateway.WorkloadClaims) { v.Scopes = []string{"audit.read", "audit.read"} },
		func(v *coregateway.WorkloadClaims) { v.TenantID, v.ActorGID = tenantID, workloadID },
	} {
		other := claims
		change(&other)
		bad := request
		bad.Token = sign(other)
		if _, err := service.Verify(ctx, bad); !errors.Is(err, scope.ErrDenied) { t.Fatal("forged binding accepted") }
	}
	bad := request
	bad.Token = "malformed"
	if _, err := service.Verify(ctx, bad); !errors.Is(err, scope.ErrDenied) { t.Fatal("malformed token accepted") }
	bad = request
	bad.Target.TenantID = callerID
	if _, err := service.Verify(ctx, bad); !errors.Is(err, scope.ErrDenied) { t.Fatal("wrong tenant accepted") }
	repo.err = errors.New("revoked or unavailable")
	if _, err := service.Verify(ctx, request); !errors.Is(err, scope.ErrDenied) { t.Fatal("revocation failed closed") }
	repo.err = nil
	repo.value.WorkloadUpdatedAt = time.Now().Add(time.Second)
	if _, err := service.Verify(ctx, request); !errors.Is(err, scope.ErrDenied) { t.Fatal("pre-reactivation token accepted") }
}

func TestTenantState(t *testing.T) {
	repo := &repository{state: authorization.TenantState{Active: false, SessionGeneration: 3}}
	service := impl.NewAuthorizationService(repo, "platform-dev")
	ctx := context.Background()
	state, err := service.TenantState(ctx, projectID, tenantID)
	if err != nil || state.Active || state.SessionGeneration != 3 { t.Fatalf("wrong state: %+v %v", state, err) }
	repo.state.SessionGeneration = 0
	if _, err := service.TenantState(ctx, projectID, tenantID); !errors.Is(err, scope.ErrDenied) { t.Fatal("invalid generation accepted") }
	if _, err := service.TenantState(ctx, "bad", tenantID); !errors.Is(err, scope.ErrDenied) { t.Fatal("invalid project accepted") }
}
