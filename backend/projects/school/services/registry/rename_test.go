package registry_test

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/generated/audit"
	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/services/delivery"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry"
	coredatabase "grovehub/libs/core/database"
	core "grovehub/libs/core/gateway"
)

func TestSchoolRename(t *testing.T) {
	dsn := os.Getenv("PROJECT_SCHOOL_TEST_DSN")
	if dsn == "" {
		t.Skip("requires disposable School database")
	}
	ctx := context.Background()
	db, err := coredatabase.OpenProjectPool(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var database string
	if err = db.QueryRowContext(ctx, "SELECT current_database()").Scan(&database); err != nil || database != "school_test" {
		t.Fatal("disposable database required")
	}
	if err = registry.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	school, other := uuid.NewString(), uuid.NewString()
	for _, id := range []string{school, other} {
		if err = registry.RegisterSchool(ctx, db, id, "school-"+id, "Original"); err != nil {
			t.Fatal(err)
		}
	}
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	key := core.WorkloadKeyConfig{KeyID: "rename-test", PublicKey: base64.RawURLEncoding.EncodeToString(pub), Issuer: "internal-tool-school-rename", Subject: "internal-tool-school-rename:" + school, Scopes: []string{"internal.school.registry.rename"}}
	raw, _ := json.Marshal([]core.WorkloadKeyConfig{key})
	mux := http.NewServeMux()
	if err = registry.RegisterHTTP(mux, db, "", string(raw)); err != nil {
		t.Fatal(err)
	}
	token, err := core.SignWorkload(key.KeyID, key.Issuer, priv, core.WorkloadClaims{Subject: key.Subject, Audience: "school", Scopes: key.Scopes, ExpiresAt: time.Now().Add(time.Minute)})
	if err != nil {
		t.Fatal(err)
	}
	correlationID := uuid.NewString()
	call := func(id, execution, expected, name string, version int64, auth bool) int {
		body, _ := json.Marshal(map[string]any{"expectedName": expected, "name": name, "expectedVersion": version, "correlationId": correlationID})
		r := httptest.NewRequest("PATCH", "/school/restricted/registry/schools/"+id+"/name", bytes.NewReader(body))
		r.Header.Set("Idempotency-Key", execution)
		if auth {
			r.Header.Set("Authorization", "Bearer "+token)
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w.Code
	}
	execution := uuid.NewString()
	for _, tc := range []struct {
		id, execution, expected, name string
		version                       int64
		auth                          bool
		want                          int
	}{
		{other, execution, "Original", "Renamed", 1, true, 403},
		{school, execution, "Original", "Renamed", 1, false, 403},
		{school, execution, "Wrong", "Renamed", 1, true, 409},
		{school, execution, "Original", "Renamed", 1, true, 204},
		{school, execution, "Original", "Renamed", 1, true, 204},
		{school, execution, "Original", "Different", 1, true, 409},
		{school, uuid.NewString(), "Original", "Different", 1, true, 409},
		{school, uuid.NewString(), "Renamed", "Final", 2, true, 204},
		{school, execution, "Original", "Renamed", 1, true, 204},
	} {
		if got := call(tc.id, tc.execution, tc.expected, tc.name, tc.version, tc.auth); got != tc.want {
			t.Fatal("unexpected rename result", tc.want, got)
		}
	}
	if got := call(school, uuid.NewString(), "Final", "Original", 3, true); got != 204 {
		t.Fatal("rename back", got)
	}
	if got := call(school, uuid.NewString(), "Original", "Stale", 1, true); got != 409 {
		t.Fatal("ABA accepted", got)
	}
	var name string
	var count int
	if err = db.QueryRowContext(ctx, "SELECT name FROM school_registry.schools WHERE id=$1", school).Scan(&name); err != nil || name != "Original" {
		t.Fatal("replay changed name", name, err)
	}
	if err = db.QueryRowContext(ctx, "SELECT count(*) FROM school_registry.rename_events WHERE school_id=$1", school).Scan(&count); err != nil || count != 3 {
		t.Fatal("wrong receipt count", count, err)
	}
	if _, err = db.ExecContext(ctx, "DELETE FROM school_registry.rename_events WHERE execution_id=$1", execution); err == nil {
		t.Fatal("immutable receipt deleted")
	}
	if _, err = db.ExecContext(ctx, "UPDATE school_registry.schools SET status='DISABLED' WHERE id=$1", school); err != nil {
		t.Fatal(err)
	}
	if got := call(school, uuid.NewString(), "Original", "Disabled change", 5, true); got != 403 {
		t.Fatal("disabled school mutated", got)
	}
	// A successful rename creates one outbox entry, not one entry per retry.
	if err = db.QueryRowContext(ctx, "SELECT count(*) FROM school_registry.rename_outbox o JOIN school_registry.rename_events e USING(execution_id) WHERE e.school_id=$1", school).Scan(&count); err != nil || count != 3 {
		t.Fatal("outbox receipt count", count, err)
	}
	health := func(configs []registry.DeliveryConfig) int {
		w := httptest.NewRecorder()
		registry.DeliveryHealth(db, configs).ServeHTTP(w, httptest.NewRequest("GET", "/health/audit-delivery", nil))
		return w.Code
	}
	if got := health(nil); got != 503 {
		t.Fatal("disabled health", got)
	}
	if got := health([]registry.DeliveryConfig{{SchoolID: other}}); got != 503 {
		t.Fatal("unconfigured backlog health", got)
	}
	calls := 0
	var first string
	worker := delivery.New(db, school, func(_ context.Context, a sdk.SchoolRenameReceipt) (int, error) {
		calls++
		if a.SchoolId != school || a.CorrelationId == nil || *a.CorrelationId != correlationID || a.PreviousVersion < 1 {
			t.Error("incorrect receipt")
		}
		if calls == 1 {
			first = a.ExecutionId
			return 0, errors.New("temporary failure")
		}
		if calls == 2 && a.ExecutionId != first {
			t.Error("retry identity changed")
		}
		return 204, nil
	})
	if worked, err := worker.Step(ctx); !worked || err != nil {
		t.Fatal("relay step", worked, err)
	}
	var attempts int
	if err = db.QueryRowContext(ctx, "SELECT attempts FROM school_registry.rename_outbox WHERE execution_id=$1 AND delivered_at IS NULL AND next_attempt_at>now()", first).Scan(&attempts); err != nil || attempts != 1 {
		t.Fatal("retry not scheduled", attempts, err)
	}
	if _, err = db.ExecContext(ctx, "UPDATE school_registry.rename_outbox SET next_attempt_at='2000-01-01' WHERE execution_id=$1", first); err != nil {
		t.Fatal(err)
	}
	for n := 0; n < 3; n++ {
		if worked, err := worker.Step(ctx); !worked || err != nil {
			t.Fatal("relay delivery", worked, err)
		}
	}
	if worked, err := worker.Step(ctx); worked || err != nil {
		t.Fatal("delivered events retried", worked, err)
	}
	if err = db.QueryRowContext(ctx, "SELECT count(*) FROM school_registry.rename_outbox o JOIN school_registry.rename_events e USING(execution_id) WHERE e.school_id=$1 AND o.delivered_at IS NOT NULL", school).Scan(&count); err != nil || count != 3 {
		t.Fatal("delivery count", count, err)
	}
}
