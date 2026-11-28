package telemetry_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coretelemetry "grovehub/libs/core/telemetry"
)

var _ = Describe("Logical service identity", func() {
	It("uses the canonical launcher edge identity", func() {
		service := coretelemetry.ResolveService("edge")

		Expect(service.Key).To(Equal("launcher-edge"))
		Expect(service.Name).To(Equal("LAUNCHER (EDGE)"))
	})

	It("uses the canonical identity service identity", func() {
		service := coretelemetry.ResolveService("identity")

		Expect(service.Key).To(Equal("identity"))
		Expect(service.Name).To(Equal("IDENTITY"))
	})

	It("carries the logical service through request context", func() {
		ctx := coretelemetry.ContextWithService(
			context.Background(),
			coretelemetry.ServiceIdentity,
		)

		service, ok := coretelemetry.ServiceFromContext(ctx)
		Expect(ok).To(BeTrue())
		Expect(service).To(Equal(coretelemetry.ServiceIdentity))
	})
})
