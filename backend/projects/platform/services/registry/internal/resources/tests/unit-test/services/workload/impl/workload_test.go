package impl_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"grovehub/libs/platform/scope"
	"grovehub/services/registry/internal/app/constants/capability"
	"grovehub/services/registry/internal/app/constants/status"
	dto "grovehub/services/registry/internal/app/dto/workload"
	"grovehub/services/registry/internal/app/entity"
	"grovehub/services/registry/internal/app/repository/constants"
	"grovehub/services/registry/internal/app/services/workload/impl"
)

const workloadID = "30000000-0000-4000-8000-000000000001"
const projectID = "10000000-0000-4000-8000-000000000001"
const tenantID = "20000000-0000-4000-8000-000000000001"
const grantID = "40000000-0000-4000-8000-000000000001"

type workloads struct {
	calls int
	value entity.Workload
	err   error
}

func (r *workloads) Find(_ context.Context, _ string) (entity.Workload, error) {
	r.calls++
	return r.value, r.err
}

func (r *workloads) Create(_ context.Context, value entity.Workload) (entity.Workload, error) {
	r.calls++
	r.value = value
	return value, r.err
}

func (r *workloads) ChangeStatus(_ context.Context, id string, version int64, next status.Status) (entity.Workload, error) {
	r.calls++
	return entity.Workload{WorkloadID: id, Version: version + 1, Status: next}, r.err
}

type keys struct {
	calls int
	value entity.WorkloadKey
	err   error
}

func (r *keys) Create(_ context.Context, value entity.WorkloadKey) error {
	r.calls++
	r.value = value
	return r.err
}

func (r *keys) Revoke(_ context.Context, _, _ string) error {
	r.calls++
	return r.err
}

type grants struct {
	calls int
	value entity.WorkloadGrant
	err   error
}

func (r *grants) Create(_ context.Context, value entity.WorkloadGrant) error {
	r.calls++
	r.value = value
	return r.err
}

func (r *grants) Revoke(_ context.Context, _, _ string) error {
	r.calls++
	return r.err
}

func authorize(t *testing.T, target scope.Target, operation string) context.Context {
	t.Helper()
	ctx, err := scope.WithVerifiedExecution(context.Background(), scope.Execution{
		Target: target, Capability: operation,
		PrincipalID: "operator-workload", PrincipalType: scope.Workload, WorkloadID: "operator-workload",
	})
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}

func TestWorkloadProvisioning(t *testing.T) {
	w, k, g := &workloads{}, &keys{}, &grants{}
	service := impl.NewWorkloadService(w, k, g)
	platform := scope.Target{Level: scope.Platform}
	project := scope.Target{Level: scope.Project, ProjectID: projectID}
	owner := projectID
	input := dto.CreateInput{ProjectID: &owner, Issuer: "synthetic", Subject: "publisher"}
	if _, err := service.Create(authorize(t, project, capability.WorkloadCreate), input); err == nil || w.calls != 0 {
		t.Fatal("project workload attempted to provision its own identity")
	}
	createCtx := authorize(t, platform, capability.WorkloadCreate)
	if _, err := service.Create(createCtx, dto.CreateInput{}); err == nil || w.calls != 0 {
		t.Fatal("invalid identity reached repository")
	}
	created, err := service.Create(createCtx, input)
	if err != nil || created.Status != status.Disabled || created.Version != 1 || created.WorkloadID == "" || w.calls != 1 {
		t.Fatalf("unsafe creation: %+v %v", created, err)
	}
	if _, err := service.Get(createCtx, workloadID); err == nil || w.calls != 1 {
		t.Fatal("create capability allowed read")
	}
	if _, err := service.Get(authorize(t, platform, capability.WorkloadRead), workloadID); err != nil {
		t.Fatal(err)
	}
	changeCtx := authorize(t, platform, capability.WorkloadStatus)
	changed, err := service.ChangeStatus(changeCtx, workloadID, dto.StatusInput{ExpectedVersion: 1, Status: status.Active})
	if err != nil || changed.Status != status.Active || changed.Version != 2 {
		t.Fatalf("invalid transition: %+v %v", changed, err)
	}
	w.err = constants.ErrConflict
	if _, err := service.ChangeStatus(changeCtx, workloadID, dto.StatusInput{ExpectedVersion: 1, Status: status.Active}); !errors.Is(err, constants.ErrConflict) {
		t.Fatalf("lost version conflict: %v", err)
	}
}

func TestKeyAndGrantProvisioning(t *testing.T) {
	w, k, g := &workloads{}, &keys{}, &grants{}
	service := impl.NewWorkloadService(w, k, g)
	platform := scope.Target{Level: scope.Platform}
	project := scope.Target{Level: scope.Project, ProjectID: projectID}
	now := time.Now().UTC()
	publicKey := bytes.Repeat([]byte{7}, 32)
	input := dto.KeyInput{KeyID: "first", PublicKey: base64.RawURLEncoding.EncodeToString(publicKey), ValidFrom: now, ValidUntil: now.Add(time.Hour)}
	if err := service.RegisterKey(authorize(t, project, capability.KeyRegister), workloadID, input); err == nil || k.calls != 0 {
		t.Fatal("project-level credential mutation allowed")
	}
	keyCtx := authorize(t, platform, capability.KeyRegister)
	invalid := input
	invalid.PublicKey = "bad"
	if err := service.RegisterKey(keyCtx, workloadID, invalid); err == nil || k.calls != 0 {
		t.Fatal("malformed key reached repository")
	}
	if err := service.RegisterKey(keyCtx, workloadID, input); err != nil || !bytes.Equal(k.value.PublicKey, publicKey) || k.value.WorkloadID != workloadID {
		t.Fatalf("key registration failed: %v", err)
	}
	if err := service.RevokeKey(keyCtx, workloadID, input.KeyID); err == nil || k.calls != 1 {
		t.Fatal("register capability allowed revocation")
	}
	if err := service.RevokeKey(authorize(t, platform, capability.KeyRevoke), workloadID, input.KeyID); err != nil {
		t.Fatal(err)
	}
	grantCtx := authorize(t, platform, capability.GrantCreate)
	grant := dto.GrantInput{Scope: scope.Tenant, TargetProjectID: projectID, TargetTenantID: tenantID, Capability: "audit.read"}
	if _, err := service.CreateGrant(authorize(t, project, capability.GrantCreate), workloadID, grant); err == nil || g.calls != 0 {
		t.Fatal("project granted itself a capability")
	}
	bad := grant
	bad.TargetTenantID = ""
	if _, err := service.CreateGrant(grantCtx, workloadID, bad); err == nil || g.calls != 0 {
		t.Fatal("tenant grant omitted its tenant")
	}
	if id, err := service.CreateGrant(grantCtx, workloadID, grant); err != nil || id == "" || g.value.TargetTenantID == nil || *g.value.TargetTenantID != tenantID {
		t.Fatalf("grant creation failed: %v", err)
	}
	if err := service.RevokeGrant(grantCtx, workloadID, grantID); err == nil || g.calls != 1 {
		t.Fatal("create capability allowed grant revocation")
	}
	if err := service.RevokeGrant(authorize(t, platform, capability.GrantRevoke), workloadID, grantID); err != nil {
		t.Fatal(err)
	}
}
