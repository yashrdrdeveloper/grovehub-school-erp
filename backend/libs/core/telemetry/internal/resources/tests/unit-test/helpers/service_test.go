package helpers_test

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ResolveService", func() {
	It("normalizes canonical service aliases and unknown names", func() {
		for _,name:=range []string{"","edge","launcher","launcher-edge","launcher (edge)","launcher(edge)"," EDGE "}{
			Expect(helpers.ResolveService(name)).To(Equal(contracts.ServiceLauncherEdge))
		}
		Expect(helpers.ResolveService(" FOUNDATION ")).To(Equal(contracts.ServiceFoundation))
		Expect(helpers.ResolveService("Identity")).To(Equal(contracts.ServiceIdentity))
		Expect(helpers.ResolveService("audit")).To(Equal(contracts.ServiceAudit))
		Expect(helpers.ResolveService(" integration ")).To(Equal(contracts.Service{Key:"integration",Name:"INTEGRATION"}))
	})
})
