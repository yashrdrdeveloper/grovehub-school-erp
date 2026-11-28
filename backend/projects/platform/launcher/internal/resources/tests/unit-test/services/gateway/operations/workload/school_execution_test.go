package workload_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"
	ginkgo "github.com/onsi/ginkgo/v2"

	. "grovehub/launcher/internal/src/services/gateway/operations/workload"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/schoolexecution"
	core "grovehub/libs/core/gateway"
)

func verifySchoolExecutionBoundary() {
	t := ginkgo.GinkgoT()
	school, other := uuid.NewString(), uuid.NewString()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	key := core.WorkloadKeyConfig{KeyID: "school-check", PublicKey: base64.RawURLEncoding.EncodeToString(pub), Issuer: schoolexecution.SchoolExecutionIssuer, Subject: schoolexecution.SchoolSubjectPrefix + school, Scopes: []string{schoolexecution.SchoolAuditCheckScope}}
	raw, _ := json.Marshal(WorkloadConfig{SchoolExecutionKeys: []core.WorkloadKeyConfig{key}})
	verifier, err := NewPlatformWorkloads(string(raw), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	token, err := core.SignWorkload(key.KeyID, key.Issuer, priv, core.WorkloadClaims{Subject: key.Subject, Audience: "platform", Scopes: key.Scopes, ExpiresAt: time.Now().Add(time.Minute)})
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		school, method, suffix, header string
		allow                          bool
	}{
		{school, "GET", "", "", true}, {other, "GET", "", "", false}, {school, "GET", "?tenant=" + other, "", false}, {school, "POST", "", "", false}, {school, "GET", "", "Cookie", false}, {school, "GET", "", "X-GH-Identity-Proof", false},
	} {
		r := httptest.NewRequest(tc.method, "/audit/restricted/schools/"+tc.school+"/activity-check"+tc.suffix, nil)
		r.Pattern = schoolexecution.SchoolAuditCheckRoute
		r.Header.Set("Authorization", "Bearer "+token)
		if tc.header != "" {
			r.Header.Set(tc.header, "forged")
		}

		_, err = verifier.VerifyWorkload(context.Background(), r)
		if (err == nil) != tc.allow {
			t.Fatal("unexpected authorization", tc, err)
		}
	}

	bad := key
	bad.Subject = schoolexecution.SchoolSubjectPrefix + other
	if _, err := schoolexecution.New([]core.WorkloadKeyConfig{key, bad}, nil); err == nil {
		t.Fatal("shared school key accepted")
	}

	if _, err := schoolexecution.New([]core.WorkloadKeyConfig{key}, []core.WorkloadKeyConfig{key}); err == nil {
		t.Fatal("shared workload accepted")
	}
}

// @Testing suite=LauncherWorkloads
func SchoolExecutionSpecs() {
	ginkgo.Describe("SchoolExecution", func() {
		ginkgo.It("School Execution Boundary", verifySchoolExecutionBoundary)
	})
}
