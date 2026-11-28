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
	"grovehub/services/registry/internal/app/constants/status"
	dto "grovehub/services/registry/internal/app/dto/project"
	"grovehub/services/registry/internal/app/repository/constants"
	"grovehub/services/registry/internal/app/services/project/impl"
)

type repository struct {
	calls int
	err   error
}

func (r *repository) Find(_ context.Context, projectID string) (entity.Project, error) {
	r.calls++
	return entity.Project{ProjectID: projectID}, r.err
}

func TestScopedRead(t *testing.T) {
	const projectID = "10000000-0000-4000-8000-000000000001"
	
	repo := &repository{}
	service := impl.NewProjectService(repo)
	ctx := context.Background()
	assertStatus := func(err error, expected int) {
		t.Helper()
		var app *coreerrors.Error
		if !errors.As(err, &app) || app.Definition().HTTPStatus() != expected {
			t.Fatalf("expected status %d, got %v", expected, err)
		}
	}
	_, err := service.Get(ctx, "bad")
	assertStatus(err, http.StatusBadRequest)
	_, err = service.Get(ctx, projectID)
	assertStatus(err, http.StatusForbidden)

	value := scope.Execution{
		Target: scope.Target{Level: scope.Tenant, ProjectID: projectID, TenantID: "20000000-0000-4000-8000-000000000001"},
		Capability: capability.ProjectRead,
		PrincipalType: scope.Workload,
		PrincipalID: "workload",
		WorkloadID: "workload",
	}
	wrong, err := scope.WithVerifiedExecution(ctx, value)
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.Get(wrong, projectID)
	assertStatus(err, http.StatusForbidden)
	if repo.calls != 0 {
		t.Fatal("unauthorized request reached repository")
	}

	value.Target = scope.Target{Level: scope.Project, ProjectID: projectID}
	verified, err := scope.WithVerifiedExecution(ctx, value)
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.Get(verified, projectID)
	if err != nil || result.ProjectID != projectID || repo.calls != 1 {
		t.Fatalf("authorized request failed: %+v %v", result, err)
	}

	repo.err = constants.ErrNotFound
	_, err = service.Get(verified, projectID)
	assertStatus(err, http.StatusNotFound)
	repo.err = errors.New("private database failure")
	_, err = service.Get(verified, projectID)
	assertStatus(err, http.StatusServiceUnavailable)
}

func (r *repository) Create(_ context.Context, value entity.Project) (entity.Project, error) {
	r.calls++
	return value, r.err
}

func (r *repository) ChangeStatus(_ context.Context, projectID string, version int64, next status.Status) (entity.Project, error) {
	r.calls++
	return entity.Project{ProjectID: projectID, Version: version + 1, Status: next}, r.err
}

func TestScopedWrites(t *testing.T) {
	const projectID = "10000000-0000-4000-8000-000000000001"
	
	repo := &repository{}
	service := impl.NewProjectService(repo)
	ctx := context.Background()
	authorized := func(target scope.Target, operation string) context.Context {
		t.Helper()
		value, err := scope.WithVerifiedExecution(ctx, scope.Execution{
			Target: target, Capability: operation,
			PrincipalType: scope.Workload, PrincipalID: "w", WorkloadID: "w",
		})
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
	input := dto.CreateInput{Code: "synthetic", Name: "Synthetic"}
	if _, err := service.Create(ctx, input); err == nil || repo.calls != 0 {
		t.Fatal("unauthorized create reached repository")
	}
	createCtx := authorized(scope.Target{Level: scope.Platform}, capability.ProjectCreate)
	if _, err := service.Create(createCtx, dto.CreateInput{}); err == nil || repo.calls != 0 {
		t.Fatal("invalid create reached repository")
	}
	created, err := service.Create(createCtx, input)
	if err != nil || created.Status != status.Disabled || created.Version != 1 || created.ProjectID == "" {
		t.Fatalf("unsafe creation: %+v %v", created, err)
	}
	change := dto.StatusInput{ExpectedVersion: 1, Status: status.Active}
	if _, err := service.ChangeStatus(createCtx, projectID, change); err == nil || repo.calls != 1 {
		t.Fatal("create grant allowed lifecycle change")
	}
	changeCtx := authorized(scope.Target{Level: scope.Project, ProjectID: projectID}, capability.ProjectStatus)
	if _, err := service.ChangeStatus(changeCtx, projectID, dto.StatusInput{}); err == nil || repo.calls != 1 {
		t.Fatal("invalid transition reached repository")
	}
	changed, err := service.ChangeStatus(changeCtx, projectID, change)
	if err != nil || changed.Status != status.Active || changed.Version != 2 {
		t.Fatalf("transition failed: %+v %v", changed, err)
	}
	repo.err = constants.ErrConflict
	_, err = service.ChangeStatus(changeCtx, projectID, change)
	var app *coreerrors.Error
	if !errors.As(err, &app) || app.Definition().HTTPStatus() != http.StatusConflict {
		t.Fatalf("wrong conflict response: %v", err)
	}
}
