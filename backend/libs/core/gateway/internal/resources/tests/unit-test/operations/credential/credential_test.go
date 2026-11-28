package credential_test

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/operations/credential"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SignWorkload and VerifyCredential", func() {
	var private ed25519.PrivateKey
	var verifier credential.Ed25519Credentials
	var claims contracts.WorkloadClaims

	BeforeEach(func() {
		public, key, err := ed25519.GenerateKey(nil)
		Expect(err).NotTo(HaveOccurred())
		private = key
		verifier = credential.Ed25519Credentials{Keys: map[string]contracts.TrustedWorkloadKey{
			"key": {PublicKey: public, Issuer: "issuer", Subject: "subject", Scopes: []string{"read"}, AllowDelegation: true},
		}}
		claims = contracts.WorkloadClaims{Subject: "subject", Audience: "platform", ExpiresAt: time.Now().Add(time.Minute), Scopes: []string{"read"}, TenantID: "tenant", ActorGID: "actor"}
	})

	It("verifies signed claims and denies altered tokens, unknown keys and canceled contexts", func() {
		token, err := credential.SignWorkload("key", "issuer", private, claims)
		Expect(err).NotTo(HaveOccurred())
		got, err := verifier.VerifyCredential(context.Background(), token)
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Subject).To(Equal(claims.Subject))
		Expect(got.Audience).To(Equal(claims.Audience))
		Expect(got.ActorGID).To(Equal(claims.ActorGID))
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err = verifier.VerifyCredential(ctx, token)
		Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		_, err = verifier.VerifyCredential(nil, token)
		Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		for _, bad := range []string{"", strings.Repeat("x", 8193), strings.Replace(token, "ghw1.", "ghw2.", 1), strings.Replace(token, ".key.", ".missing.", 1), token[:strings.LastIndex(token, ".")+1]+"!", token[:strings.LastIndex(token, ".")+1]+"YQ"} {
			_, err = verifier.VerifyCredential(context.Background(), bad)
			Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		}
	})

	It("rejects invalid signing inputs and oversized credentials", func() {
		for _, keyID := range []string{"", "a.b", strings.Repeat("x", 129)} {
			_, err := credential.SignWorkload(keyID, "issuer", private, claims)
			Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		}
		_, err := credential.SignWorkload("key", "", private, claims)
		Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		_, err = credential.SignWorkload("key", "issuer", nil, claims)
		Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		for _, mutate := range []func(*contracts.WorkloadClaims){
			func(v *contracts.WorkloadClaims) { v.Subject = "" },
			func(v *contracts.WorkloadClaims) { v.Audience = "" },
			func(v *contracts.WorkloadClaims) { v.ExpiresAt = time.Now().Add(-time.Minute) },
			func(v *contracts.WorkloadClaims) { v.ExpiresAt = time.Now().Add(3*time.Minute) },
			func(v *contracts.WorkloadClaims) { v.ActorGID = "" },
			func(v *contracts.WorkloadClaims) { v.Scopes = []string{strings.Repeat("x", 9000)} },
		} {
			value := claims
			mutate(&value)
			_, err := credential.SignWorkload("key", "issuer", private, value)
			Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		}
	})

	It("denies signed payloads that violate issuer, subject, time, scope or delegation policy", func() {
		now := time.Now().Unix()
		makeToken := func(payload string) string {
			body := "ghw1.key." + base64.RawURLEncoding.EncodeToString([]byte(payload))
			return body + "." + base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, []byte(body)))
		}
		valid := map[string]any{"iss":"issuer","sub":"subject","aud":"platform","iat":now,"exp":now+60,"scopes":[]string{"read"}}
		for field, values := range map[string][]any{
			"iss":{"other"}, "sub":{"other"}, "aud":{""},
			"iat":{int64(0),now+1}, "exp":{now,now+121},
			"scopes":{[]string{"write"},[]string{""}}, "tenant":{"tenant"},
		} {
			for _, value := range values {
				payload := make(map[string]any)
				for k, v := range valid { payload[k] = v }
				payload[field] = value
				raw, err := json.Marshal(payload)
				Expect(err).NotTo(HaveOccurred())
				_, err = verifier.VerifyCredential(context.Background(), makeToken(string(raw)))
				Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
			}
		}
		_, err := verifier.VerifyCredential(context.Background(), makeToken("not-json"))
		Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		token, err := credential.SignWorkload("key", "issuer", private, claims)
		Expect(err).NotTo(HaveOccurred())
		key := verifier.Keys["key"]
		key.AllowDelegation = false
		verifier.Keys["key"] = key
		_, err = verifier.VerifyCredential(context.Background(), token)
		Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
	})
})
