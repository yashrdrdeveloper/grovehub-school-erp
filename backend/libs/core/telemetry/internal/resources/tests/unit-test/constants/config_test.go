package constants_test

import (
	"grovehub/libs/core/telemetry/internal/src/constants"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Telemetry defaults", func() {
	It("preserves deployment, sampling and export defaults", func() {
		Expect(constants.ServiceNamespace).To(Equal("grovehub"))
		Expect(constants.DefaultEnvironment).To(Equal("local"))
		Expect(constants.DefaultServiceVersion).To(Equal("dev"))
		Expect(constants.DefaultTraceSampleRatio).To(Equal(1.0))
		Expect(constants.DefaultMetricExportInterval).To(Equal(15 * time.Second))
		Expect(constants.MinMetricExportInterval).To(Equal(5 * time.Second))
		Expect(constants.MetricCardinalityLimit).To(Equal(1000))
	})
})
