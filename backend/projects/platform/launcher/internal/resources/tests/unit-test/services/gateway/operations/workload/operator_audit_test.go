package workload_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"

	core "grovehub/libs/core/gateway"
	. "grovehub/launcher/internal/src/services/gateway/operations/workload"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/internalcatalog"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/operatoraudit"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/operatorauditread"
)

func verifyOperatorAuditBoundary() {
	t := ginkgo.GinkgoT()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	key := core.WorkloadKeyConfig{KeyID: "operator-audit-test", PublicKey: base64.RawURLEncoding.EncodeToString(pub), Issuer: "internal-tool-audit", Subject: "internal-tool-audit", Scopes: []string{operatoraudit.OperatorAuditScope}}
	raw, err := json.Marshal(WorkloadConfig{OperatorAuditKeys: []core.WorkloadKeyConfig{key}})
	if err != nil {
		t.Fatal(err)
	}

	verifier, err := NewPlatformWorkloads(string(raw), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	token, err := core.SignWorkload(key.KeyID, key.Issuer, priv, core.WorkloadClaims{Subject: key.Subject, Audience: "platform", Scopes: key.Scopes, ExpiresAt: time.Now().Add(time.Minute)})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		pattern, path, header string
		allowed               bool
	}{
		{operatoraudit.OperatorAuditRoute, "/audit/restricted/operator-activities", "", true},
		{operatoraudit.OperatorAuditRoute, "/audit/restricted/operator-activities?tenant=x", "", false},
		{operatoraudit.OperatorAuditRoute, "/audit/restricted/operator-activities", "Cookie", false},
		{operatoraudit.OperatorAuditRoute, "/audit/restricted/operator-activities", "X-GH-Identity-Proof", false},
		{"POST /audit/restricted/search", "/audit/restricted/search", "", false},
		{internalcatalog.InternalCatalogRoute, "/integrations/restricted/providers", "", false},
	} {
		r := httptest.NewRequest("POST", tc.path, nil)
		r.Pattern = tc.pattern
		r.Header.Set("Authorization", "Bearer "+token)
		if tc.header != "" {
			r.Header.Set(tc.header, "forged")
		}

		_, err := verifier.VerifyWorkload(context.Background(), r)
		if (err == nil) != tc.allowed {
			t.Errorf("%s %s: %v", tc.pattern, tc.header, err)
		}
	}

	bad := key
	bad.Scopes = []string{operatoraudit.OperatorAuditScope, "audit.search"}
	if _, err := operatoraudit.New([]core.WorkloadKeyConfig{bad}, nil); err == nil {
		t.Fatal("overbroad key accepted")
	}

	if _, err := operatoraudit.New([]core.WorkloadKeyConfig{key}, []core.WorkloadKeyConfig{key}); err == nil {
		t.Fatal("shared key accepted")
	}
}

func verifyOperatorAuditReadIsolation() {
	t := ginkgo.GinkgoT()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	key := core.WorkloadKeyConfig{KeyID: "read-test", PublicKey: base64.RawURLEncoding.EncodeToString(pub), Issuer: "internal-tool-audit-reader", Subject: "internal-tool-audit-reader", Scopes: []string{operatorauditread.OperatorAuditReadScope}}
	raw, _ := json.Marshal(WorkloadConfig{OperatorAuditReadKeys: []core.WorkloadKeyConfig{key}})
	verifier, err := NewPlatformWorkloads(string(raw), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	token, err := core.SignWorkload(key.KeyID, key.Issuer, priv, core.WorkloadClaims{Subject: key.Subject, Audience: "platform", Scopes: key.Scopes, ExpiresAt: time.Now().Add(time.Minute)})
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{"/audit/restricted/operator-activities/search", "/audit/restricted/operator-activities", "/audit/restricted/search"} {
		r := httptest.NewRequest("POST", path, nil)
		r.Pattern = "POST " + path
		r.Header.Set("Authorization", "Bearer "+token)
		_, err = verifier.VerifyWorkload(context.Background(), r)
		if (err == nil) != (r.Pattern == operatorauditread.OperatorAuditReadRoute) {
			t.Fatal("wrong grant", path, err)
		}
	}

	if _, err := operatorauditread.New([]core.WorkloadKeyConfig{key}, []core.WorkloadKeyConfig{key}); err == nil {
		t.Fatal("shared key accepted")
	}
}

// @Testing suite=LauncherWorkloads
func OperatorAuditSpecs() {
	ginkgo.Describe("OperatorAudit", func() {
		ginkgo.It("Operator Audit Boundary", verifyOperatorAuditBoundary)
		ginkgo.It("Operator Audit Read Isolation", verifyOperatorAuditReadIsolation)
	})
}
