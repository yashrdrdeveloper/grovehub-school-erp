package impl_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/platform/scope"
	"grovehub/services/registry/internal/app/constants/capability"
	"grovehub/services/registry/internal/app/entity"
	"grovehub/services/registry/internal/app/repository/constants"
	"grovehub/services/registry/internal/app/services/tenant/impl"
)

type repository struct {
	calls int
	err   error
}

func (r *repository) Find(_ context.Context, projectID, tenantID string) (entity.Tenant, error) {
	r.calls++
	return entity.Tenant{ProjectID: projectID, TenantID: tenantID}, r.err
}

func TestScopedRead(t *testing.T) {
	const projectID = "10000000-0000-4000-8000-000000000001"
	const tenantID = "20000000-0000-4000-8000-000000000001"
	repo := &repository{}
	service := impl.NewTenantService(repo)
	ctx := context.Background()
	assertStatus := func(err error, expected int) {
		t.Helper()
		var app *coreerrors.Error
		if !errors.As(err, &app) || app.Definition().HTTPStatus() != expected {
			t.Fatalf("expected status %d, got %v", expected, err)
		}
	}
	_, err := service.Get(ctx, projectID, "bad")
	assertStatus(err, http.StatusBadRequest)
	_, err = service.Get(ctx, projectID, tenantID)
	assertStatus(err, http.StatusForbidden)

	value := scope.Execution{
		Target: scope.Target{Level: scope.Project, ProjectID: projectID},
		Capability: capability.TenantRead,
		PrincipalType: scope.Workload,
		PrincipalID: "workload",
		WorkloadID: "workload",
	}
	wrong, err := scope.WithVerifiedExecution(ctx, value)
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.Get(wrong, projectID, tenantID)
	assertStatus(err, http.StatusForbidden)
	if repo.calls != 0 {
		t.Fatal("unauthorized request reached repository")
	}

	value.Target = scope.Target{Level: scope.Tenant, ProjectID: projectID, TenantID: tenantID}
	verified, err := scope.WithVerifiedExecution(ctx, value)
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.Get(verified, projectID, tenantID)
	if err != nil || result.ProjectID != projectID || repo.calls != 1 {
		t.Fatalf("authorized request failed: %+v %v", result, err)
	}

	repo.err = constants.ErrNotFound
	_, err = service.Get(verified, projectID, tenantID)
	assertStatus(err, http.StatusNotFound)
	repo.err = errors.New("private database failure")
	_, err = service.Get(verified, projectID, tenantID)
	assertStatus(err, http.StatusServiceUnavailable)
}
