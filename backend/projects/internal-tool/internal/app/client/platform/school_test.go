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

	"github.com/google/uuid"
	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
	core "grovehub/libs/core/gateway"
)

func TestSchoolCheckAdapter(t *testing.T) {
	school, template, operator := uuid.NewString(), uuid.NewString(), uuid.NewString()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	creds, err := core.NewEd25519Credentials([]core.WorkloadKeyConfig{{KeyID: "school-test", PublicKey: base64.RawURLEncoding.EncodeToString(pub), Issuer: schoolIssuer, Subject: schoolIssuer + ":" + school, Scopes: []string{schoolCheckScope}}})
	if err != nil {
		t.Fatal(err)
	}
	path := "/audit/restricted/schools/" + school + "/activity-check"
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		claims, err := creds.VerifyCredential(r.Context(), strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		if err != nil || claims.Subject != schoolIssuer+":"+school || len(claims.Scopes) != 1 || claims.Scopes[0] != schoolCheckScope || r.Method != "GET" || r.URL.Path != path || r.Header.Get("Cookie") != "" {
			t.Error("wrong outbound school boundary", err)
		}
		w.WriteHeader(204)
	}))
	defer server.Close()
	c := SchoolConfig{EnvironmentCode: "dev", TemplateID: template, SchoolID: school, Origin: server.URL, KeyID: "school-test", PrivateKey: base64.RawURLEncoding.EncodeToString(priv), OperatorIDs: operator}
	bindings, err := SchoolBindings(c)
	if err != nil {
		t.Fatal(err)
	}
	binding := bindings[workflow.SchoolBindingKey{EnvironmentCode: "dev", TemplateID: template, SchoolID: school}]
	call := workflow.Invocation{SchoolID: &school, Method: "GET", Path: path, Input: json.RawMessage(`{"schoolId":"` + school + `"}`)}
	if status, err := binding.Adapter(context.Background(), call); err != nil || status != 204 {
		t.Fatal("valid operation failed", status, err)
	}
	other := uuid.NewString()
	call.SchoolID = &other
	if _, err = binding.Adapter(context.Background(), call); err == nil || calls != 1 {
		t.Fatal("cross-school dispatch accepted")
	}
	call.SchoolID = &school
	call.Input = json.RawMessage(`{"schoolId":"` + other + `"}`)
	if binding.Validate(call) {
		t.Fatal("mismatched path input accepted")
	}
	c.OperatorIDs = operator + "," + operator
	if _, err = SchoolBindings(c); err == nil {
		t.Fatal("duplicate operator accepted")
	}
}
