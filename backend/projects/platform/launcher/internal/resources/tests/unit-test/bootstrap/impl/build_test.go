package impl_test

import (
	ginkgo "github.com/onsi/ginkgo/v2"
	"grovehub/launcher/internal/src/bootstrap/impl"
)

func verifyBuildRejectsMissingTenantConfiguration() {
	t := ginkgo.GinkgoT()
	t.Setenv("PLATFORM_TENANT_ACCESS_REQUIRED", "true")
	t.Setenv("PLATFORM_TENANT_ACCESS_ORIGIN", "")
	t.Setenv("PLATFORM_TENANT_ACCESS_KEY_ID", "")
	t.Setenv("PLATFORM_TENANT_ACCESS_PRIVATE_KEY", "")
	runtime, err := impl.New(nil).Build()
	if runtime != nil || err == nil {
		t.Fatal("missing required tenant verifier accepted")
	}
}

func verifyBuildRejectsMissingDatabase() {
	t := ginkgo.GinkgoT()
	t.Setenv("PLATFORM_TENANT_ACCESS_REQUIRED", "false")
	t.Setenv("PLATFORM_TENANT_ACCESS_ORIGIN", "")
	t.Setenv("PLATFORM_TENANT_ACCESS_KEY_ID", "")
	t.Setenv("PLATFORM_TENANT_ACCESS_PRIVATE_KEY", "")
	t.Setenv("DATABASE_URL", "")
	runtime, err := impl.New(nil).Build()
	if runtime != nil || err == nil {
		t.Fatal("missing database URL accepted")
	}
}

// @Testing suite=LauncherBootstrapImplementation
func BuildSpecs() {
	ginkgo.Describe("Build", func() {
		ginkgo.It("Build Rejects Missing Tenant Configuration", verifyBuildRejectsMissingTenantConfiguration)
		ginkgo.It("Build Rejects Missing Database", verifyBuildRejectsMissingDatabase)
	})
}
