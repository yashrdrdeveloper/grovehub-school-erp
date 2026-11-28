package config_test

import (
	"crypto/ed25519"
	"encoding/base64"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewEd25519Credentials", func() {
	var valid contracts.WorkloadKeyConfig

	BeforeEach(func() {
		key, _, err := ed25519.GenerateKey(nil)
		Expect(err).NotTo(HaveOccurred())
		valid = contracts.WorkloadKeyConfig{
			KeyID: "key", PublicKey: base64.RawURLEncoding.EncodeToString(key),
			Issuer: "issuer", Subject: "subject", Scopes: []string{"read"}, AllowDelegation: true,
		}
	})

	It("constructs trusted keys and copies scope configuration", func() {
		credentials, err := config.NewEd25519Credentials([]contracts.WorkloadKeyConfig{valid})
		Expect(err).NotTo(HaveOccurred())
		key := credentials.Keys["key"]
		Expect(key.Issuer).To(Equal("issuer"))
		Expect(key.Subject).To(Equal("subject"))
		Expect(key.PublicKey).To(HaveLen(ed25519.PublicKeySize))
		Expect(key.AllowDelegation).To(BeTrue())
		valid.Scopes[0] = "modified"
		Expect(key.Scopes).To(Equal([]string{"read"}))
	})

	It("rejects missing configuration and duplicate IDs", func() {
		_, err := config.NewEd25519Credentials(nil)
		Expect(err).To(HaveOccurred())
		_, err = config.NewEd25519Credentials([]contracts.WorkloadKeyConfig{valid, valid})
		Expect(err).To(HaveOccurred())
	})

	It("rejects malformed key identity, key material and scope policies", func() {
		for _, mutate := range []func(*contracts.WorkloadKeyConfig){
			func(v *contracts.WorkloadKeyConfig) { v.KeyID = "" },
			func(v *contracts.WorkloadKeyConfig) { v.KeyID = "a.b" },
			func(v *contracts.WorkloadKeyConfig) { v.Issuer = " " },
			func(v *contracts.WorkloadKeyConfig) { v.Subject = "" },
			func(v *contracts.WorkloadKeyConfig) { v.Scopes = nil },
			func(v *contracts.WorkloadKeyConfig) { v.PublicKey = "!" },
			func(v *contracts.WorkloadKeyConfig) { v.PublicKey = "YQ" },
			func(v *contracts.WorkloadKeyConfig) { v.Scopes = []string{" "} },
			func(v *contracts.WorkloadKeyConfig) { v.Scopes = []string{"read", "read"} },
		} {
			value := valid
			mutate(&value)
			_, err := config.NewEd25519Credentials([]contracts.WorkloadKeyConfig{value})
			Expect(err).To(HaveOccurred())
		}
	})
})
