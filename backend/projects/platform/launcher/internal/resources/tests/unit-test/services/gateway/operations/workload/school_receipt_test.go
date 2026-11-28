package workload_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	ginkgo "github.com/onsi/ginkgo/v2"
	core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	. "grovehub/launcher/internal/src/services/gateway/operations/workload"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/schoolreceipt"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/schoolreceiptread"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/schoolstatusreceipt"
	"net/http/httptest"
	"time"
)

func verifySchoolReceiptWorkload() {
	t := ginkgo.GinkgoT()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	k := core.WorkloadKeyConfig{KeyID: "receipt-test", PublicKey: base64.RawURLEncoding.EncodeToString(pub), Issuer: schoolreceipt.ReceiptIssuer, Subject: schoolreceipt.ReceiptIssuer + ":" + uuid.NewString(), Scopes: []string{schoolreceipt.ReceiptScope}}
	raw, _ := json.Marshal(WorkloadConfig{SchoolReceiptKeys: []core.WorkloadKeyConfig{k}})
	v, err := NewPlatformWorkloads(string(raw), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	token, err := core.SignWorkload(k.KeyID, k.Issuer, priv, core.WorkloadClaims{Subject: k.Subject, Audience: "platform", Scopes: k.Scopes, ExpiresAt: time.Now().Add(time.Minute)})
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{"/audit/restricted/school-rename-receipts", "/audit/restricted/operator-activities", "/audit/restricted/search"} {
		r := httptest.NewRequest("POST", path, nil)
		r.Pattern = "POST " + path
		r.Header.Set("Authorization", "Bearer "+token)
		_, err = v.VerifyWorkload(context.Background(), r)
		if (err == nil) != (r.Pattern == schoolreceipt.ReceiptRoute) {
			t.Fatal("wrong grant", path, err)
		}
	}

	if _, err := schoolreceipt.New([]core.WorkloadKeyConfig{k}, []core.WorkloadKeyConfig{k}); err == nil {
		t.Fatal("shared credential accepted")
	}
}

func verifySchoolReceiptReaderWorkload() {
	t := ginkgo.GinkgoT()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	k := core.WorkloadKeyConfig{KeyID: "receipt-reader-test", PublicKey: base64.RawURLEncoding.EncodeToString(pub), Issuer: "internal-tool-school-receipt-reader", Subject: "internal-tool-school-receipt-reader", Scopes: []string{schoolreceiptread.ReceiptReadScope}}
	raw, _ := json.Marshal(WorkloadConfig{SchoolReceiptReadKeys: []core.WorkloadKeyConfig{k}})
	v, err := NewPlatformWorkloads(string(raw), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	token, err := core.SignWorkload(k.KeyID, k.Issuer, priv, core.WorkloadClaims{Subject: k.Subject, Audience: "platform", Scopes: k.Scopes, ExpiresAt: time.Now().Add(time.Minute)})
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{"/audit/restricted/school-mutation-receipts/search", "/audit/restricted/school-rename-receipts", "/audit/restricted/operator-activities/search", "/audit/restricted/search"} {
		r := httptest.NewRequest("POST", path, nil)
		r.Pattern = "POST " + path
		r.Header.Set("Authorization", "Bearer "+token)
		_, err = v.VerifyWorkload(context.Background(), r)
		if (err == nil) != (r.Pattern == schoolreceiptread.ReceiptReadRoute) {
			t.Fatal("wrong reader grant", path, err)
		}
	}

	if _, err := schoolreceiptread.New([]core.WorkloadKeyConfig{k}, []core.WorkloadKeyConfig{k}); err == nil {
		t.Fatal("reused key accepted")
	}
}

func verifyStatusReceiptWorkload() {
	t := ginkgo.GinkgoT()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	k := core.WorkloadKeyConfig{KeyID: "status-receipt-test", PublicKey: base64.RawURLEncoding.EncodeToString(pub), Issuer: schoolstatusreceipt.StatusReceiptIssuer, Subject: schoolstatusreceipt.StatusReceiptIssuer + ":" + uuid.NewString(), Scopes: []string{schoolstatusreceipt.StatusReceiptScope}}
	raw, _ := json.Marshal(WorkloadConfig{SchoolStatusReceiptKeys: []core.WorkloadKeyConfig{k}})
	v, err := NewPlatformWorkloads(string(raw), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	token, err := core.SignWorkload(k.KeyID, k.Issuer, priv, core.WorkloadClaims{Subject: k.Subject, Audience: "platform", Scopes: k.Scopes, ExpiresAt: time.Now().Add(time.Minute)})
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{"/audit/restricted/school-status-receipts", "/audit/restricted/school-rename-receipts", "/audit/restricted/school-mutation-receipts/search"} {
		r := httptest.NewRequest("POST", path, nil)
		r.Pattern = "POST " + path
		r.Header.Set("Authorization", "Bearer "+token)
		_, err = v.VerifyWorkload(context.Background(), r)
		if (err == nil) != (r.Pattern == schoolstatusreceipt.StatusReceiptRoute) {
			t.Fatal("wrong status publisher grant", path, err)
		}
	}
}

// @Testing suite=LauncherWorkloads
func SchoolReceiptSpecs() {
	ginkgo.Describe("SchoolReceipt", func() {
		ginkgo.It("School Receipt Workload", verifySchoolReceiptWorkload)
		ginkgo.It("School Receipt Reader Workload", verifySchoolReceiptReaderWorkload)
		ginkgo.It("Status Receipt Workload", verifyStatusReceiptWorkload)
	})
}
