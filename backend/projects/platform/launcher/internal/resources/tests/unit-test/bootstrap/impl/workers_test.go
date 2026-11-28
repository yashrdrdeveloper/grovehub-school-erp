package impl_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"grovehub/launcher/internal/src/bootstrap"
	"grovehub/launcher/internal/src/bootstrap/impl"
	"grovehub/launcher/internal/src/config"
)

// @Testing suite=LauncherBootstrapImplementation
func WorkerSpecs() {
	Describe("Worker ownership", func() {
		It("starts the tenant reconciler only with a validated policy and stops it before dependencies after server failure", func() {
			_, key, err := ed25519.GenerateKey(rand.Reader)
			Expect(err).NotTo(HaveOccurred())
			fixture := newCompositionFixture("")
			cfg := config.Config{
				Address: "127.0.0.1:0", TenantRequired: true,
				TenantOrigin: "https://school.example", TenantKeyID: "tenant",
				TenantPrivateKey: base64.RawURLEncoding.EncodeToString(key),
			}
			failure := errors.New("server failed")
			fixture.server.EXPECT().Run(gomock.Any()).Return(failure)
			Expect(bootstrap.Run(impl.NewWithDependencies(compositionLogger(), cfg, fixture.dependencies))).To(MatchError(failure))
			Expect(fixture.started).To(ContainElement("tenant worker"))
			Expect(fixture.closed).To(Equal([]string{"tenant worker", "audit", "audit database", "identity", "integration", "foundation", "infrastructure"}))
		})
		It("rejects invalid tenant credentials before opening the database", func() {
			cfg := config.Config{TenantRequired: true, TenantOrigin: "http://invalid.example", TenantKeyID: "key", TenantPrivateKey: "invalid"}
			fixture := newCompositionFixture("")
			runtime, err := impl.NewWithDependencies(compositionLogger(), cfg, fixture.dependencies).Build()
			Expect(runtime).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("invalid tenant access origin")))
			Expect(fixture.started).To(BeEmpty())
		})
	})
}
