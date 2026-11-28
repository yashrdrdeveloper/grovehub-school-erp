package contracts_test

import (
	"context"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service contracts", func() {
	It("preserves logical service identities and shutdown signature", func() {
		Expect(contracts.ServiceLauncherEdge).To(Equal(contracts.Service{Key:"launcher-edge",Name:"LAUNCHER (EDGE)"}))
		Expect(contracts.ServiceFoundation).To(Equal(contracts.Service{Key:"foundation",Name:"FOUNDATION"}))
		Expect(contracts.ServiceIdentity).To(Equal(contracts.Service{Key:"identity",Name:"IDENTITY"}))
		Expect(contracts.ServiceAudit).To(Equal(contracts.Service{Key:"audit",Name:"AUDIT"}))
		var shutdown contracts.Shutdown = func(context.Context) error { return nil }
		Expect(shutdown(context.Background())).To(Succeed())
	})
})
