package schema_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"grovehub/services/registry/internal/app/repository/constants"
	projectrepo "grovehub/services/registry/internal/app/repository/project/impl"
	tenantrepo "grovehub/services/registry/internal/app/repository/tenant/impl"
	workloadrepo "grovehub/services/registry/internal/app/repository/workload/impl"

	"grovehub/services/registry/internal/resources/migrations"
)

func TestRegistryConstraints(t *testing.T) {
	dsn := os.Getenv("REGISTRY_TEST_DSN")
	if dsn == "" {
		t.Skip("REGISTRY_TEST_DSN requires a disposable database")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	var name string
	if err := db.QueryRowContext(ctx, "SELECT current_database()").Scan(&name); err != nil || name != "groovehub_registry_test" {
		t.Fatalf("disposable registry database required: %q %v", name, err)
	}
	if err := migrations.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	if err := migrations.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			t.Fatal(err)
		}
	}
	reject := func(query string, args ...any) {
		t.Helper()
		exec("SAVEPOINT invalid_insert")
		if _, err := tx.ExecContext(ctx, query, args...); err == nil {
			t.Fatal("accepted invalid registry record")
		}
		exec("ROLLBACK TO SAVEPOINT invalid_insert")
		exec("RELEASE SAVEPOINT invalid_insert")
	}

	const p1 = "10000000-0000-4000-8000-000000000001"
	const p2 = "10000000-0000-4000-8000-000000000002"
	const tenant = "20000000-0000-4000-8000-000000000001"
	const workload = "30000000-0000-4000-8000-000000000001"
	const grant1 = "40000000-0000-4000-8000-000000000001"
	const grant2 = "40000000-0000-4000-8000-000000000002"
	exec("INSERT INTO registry.projects(project_id, code, name) VALUES ($1,'synthetic-a','Synthetic A'),($2,'synthetic-b','Synthetic B')", p1, p2)
	exec("INSERT INTO registry.tenants(tenant_id,project_id,external_reference) VALUES ($1,$2,'tenant-one')", tenant, p1)
	exec("INSERT INTO registry.workloads(workload_id,project_id,issuer,subject) VALUES ($1,$2,'fixture','executor')", workload, p2)
	exec("INSERT INTO registry.capabilities(code,scope) VALUES ('audit.read','TENANT'),('audit.read','PROJECT')")
	const insertGrant = "INSERT INTO registry.workload_grants(grant_id,workload_id,target_project_id,target_tenant_id,scope,capability) VALUES ($1,$2,$3,$4,$5,$6)"
	exec(insertGrant, grant1, workload, p1, tenant, "TENANT", "audit.read")
	reject(insertGrant, grant2, workload, p2, tenant, "TENANT", "audit.read")
	reject(insertGrant, grant2, workload, p1, nil, "TENANT", "audit.read")
	reject(insertGrant, grant2, workload, p1, tenant, "PROJECT", "audit.read")
	reject(insertGrant, grant2, workload, p1, tenant, "TENANT", "unregistered")
	reject(insertGrant, grant2, workload, p1, tenant, "TENANT", "audit.read")
	exec(insertGrant, grant2, workload, p1, nil, "PROJECT", "audit.read")

	var status string
	var generation int64
	if err := tx.QueryRowContext(ctx, "SELECT status, session_generation FROM registry.tenants WHERE tenant_id=$1", tenant).Scan(&status, &generation); err != nil || status != "DISABLED" || generation != 1 {
		t.Fatalf("unsafe tenant defaults: %q %d %v", status, generation, err)
	}
	// Run GORM repository reads inside the same transaction as the fixtures.
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: tx}), &gorm.Config{
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	project, err := projectrepo.NewProjectRepository(gormDB).Find(ctx, p1)
	if err != nil || project.ProjectID != p1 || project.Code != "synthetic-a" {
		t.Fatalf("project lookup failed: %+v %v", project, err)
	}
	tenants := tenantrepo.NewTenantRepository(gormDB)
	record, err := tenants.Find(ctx, p1, tenant)
	if err != nil || record.ProjectID != p1 || record.TenantID != tenant {
		t.Fatalf("tenant lookup failed: %+v %v", record, err)
	}
	if _, err := tenants.Find(ctx, p2, tenant); !errors.Is(err, constants.ErrNotFound) {
		t.Fatalf("cross-project tenant lookup was not rejected: %v", err)
	}
	if _, err := tenants.Find(ctx, p1, "20000000-0000-4000-8000-000000000099"); !errors.Is(err, constants.ErrNotFound) {
		t.Fatalf("missing tenant lookup returned %v", err)
	}
	caller, err := workloadrepo.NewWorkloadRepository(gormDB).Find(ctx, workload)
	if err != nil || caller.ProjectID == nil || *caller.ProjectID != p2 || caller.WorkloadID != workload {
		t.Fatalf("workload lookup failed: %+v %v", caller, err)
	}

}
