package registry_test

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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

func TestRegistry(t *testing.T) {
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
	if err = registry.Migrate(ctx, db); err != nil {
		t.Fatal(err)
	}
	if err = registry.Check(ctx, db); err != nil {
		t.Fatal(err)
	}
	id, other, missing := uuid.NewString(), uuid.NewString(), uuid.NewString()
	code := "school-" + id
	if err = registry.RegisterSchool(ctx, db, id, code, "Test School"); err != nil {
		t.Fatal(err)
	}
	if err = registry.RegisterSchool(ctx, db, id, code, "Test School"); err != nil {
		t.Fatal("idempotent registration", err)
	}
	if err = registry.RegisterSchool(ctx, db, id, code, "Changed Name"); err == nil {
		t.Fatal("unexpected overwrite")
	}
	if err = registry.RegisterSchool(ctx, db, other, code, "Other"); err == nil {
		t.Fatal("duplicate code")
	}
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	key := core.WorkloadKeyConfig{KeyID: "directory-test", PublicKey: base64.RawURLEncoding.EncodeToString(pub), Issuer: "internal-tool-school-directory", Subject: "internal-tool-school-directory:" + uuid.NewString(), Scopes: []string{"internal.school.registry.read"}}
	raw, _ := json.Marshal([]map[string]any{{"key": key, "schoolIds": []string{id, missing}}})
	mux := http.NewServeMux()
	if err = registry.RegisterHTTP(mux, db, string(raw)); err != nil {
		t.Fatal(err)
	}
	token, err := core.SignWorkload(key.KeyID, key.Issuer, priv, core.WorkloadClaims{Subject: key.Subject, Audience: "school", Scopes: key.Scopes, ExpiresAt: time.Now().Add(time.Minute)})
	if err != nil {
		t.Fatal(err)
	}
	lookup := func(ids []string, auth bool, header string) *httptest.ResponseRecorder {
		body, _ := json.Marshal(map[string]any{"schoolIds": ids})
		r := httptest.NewRequest("POST", "/school/restricted/registry/lookup", bytes.NewReader(body))
		if auth {
			r.Header.Set("Authorization", "Bearer "+token)
		}
		if header != "" {
			r.Header.Set(header, "forged")
		}
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w
	}
	w := lookup([]string{id, missing}, true, "")
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	var response struct {
		Records []struct{ ID, Code, Name, Status string }
	}
	if err = json.Unmarshal(w.Body.Bytes(), &response); err != nil || len(response.Records) != 1 || response.Records[0].ID != id || response.Records[0].Name != "Test School" || response.Records[0].Status != "ACTIVE" {
		t.Fatal("incorrect lookup", response, err)
	}
	for _, tc := range []struct {
		ids    []string
		auth   bool
		header string
		status int
	}{
		{[]string{other}, true, "", 403}, {[]string{id, other}, true, "", 403}, {[]string{id, id}, true, "", 400}, {[]string{id}, false, "", 403}, {[]string{id}, true, "Cookie", 403}, {[]string{id}, true, "X-GH-Identity-Proof", 403},
	} {
		if got := lookup(tc.ids, tc.auth, tc.header).Code; got != tc.status {
			t.Fatal("wrong lookup status", tc, got)
		}
	}
}
