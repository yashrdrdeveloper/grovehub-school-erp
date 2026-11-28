package platform

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
	core "grovehub/libs/core/gateway"
)

func TestCatalogAdapter(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	creds, err := core.NewEd25519Credentials([]core.WorkloadKeyConfig{{KeyID: "internal-test", PublicKey: base64.RawURLEncoding.EncodeToString(pub), Issuer: "internal-tool", Subject: "internal-tool", Scopes: []string{catalogScope}}})
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != "GET" || r.URL.Path != catalogPath || r.URL.RawQuery != "" || r.Header.Get("Cookie") != "" || r.Header.Get("X-GH-Identity-Proof") != "" {
			t.Error("unexpected outbound request")
		}
		claims, err := creds.VerifyCredential(r.Context(), strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		if err != nil || claims.Subject != "internal-tool" || len(claims.Scopes) != 1 || claims.Scopes[0] != catalogScope {
			t.Error("incorrect workload credential", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":[{"id":"GMAIL","displayName":"Gmail","capabilities":[],"authType":"OAUTH2"}]}`))
	}))
	defer server.Close()
	id := "6cb6dcb6-f736-4c3f-90a0-131c16efdb57"
	bindings, err := Bindings(Config{EnvironmentCode: "dev", TemplateID: id, Origin: server.URL, KeyID: "internal-test", PrivateKey: base64.RawURLEncoding.EncodeToString(priv)})
	if err != nil {
		t.Fatal(err)
	}
	adapter := bindings[workflow.BindingKey{EnvironmentCode: "dev", TemplateID: id}]
	if adapter == nil {
		t.Fatal("missing binding")
	}
	call := workflow.Invocation{Method: "GET", Path: catalogPath, Input: json.RawMessage("{}")}
	status, err := adapter(context.Background(), call)
	if err != nil || status != 200 {
		t.Fatal("catalog request failed", status, err)
	}
	school := "school"
	call.SchoolID = &school
	if _, err := adapter(context.Background(), call); err == nil {
		t.Fatal("school scope accepted")
	}
	call.SchoolID = nil
	call.Path = "/identity/restricted/accounts"
	if _, err := adapter(context.Background(), call); err == nil || calls != 1 {
		t.Fatal("unbound request dispatched")
	}
}
