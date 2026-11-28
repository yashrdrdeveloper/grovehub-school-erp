package telemetry_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coretelemetry "grovehub/libs/core/telemetry"
)

var _ = Describe("Telemetry configuration facade",func(){
	It("exposes former configuration validation APIs through the root",func(){
		ratio,err:=coretelemetry.TraceSampleRatio("")
		Expect(err).NotTo(HaveOccurred())
		Expect(ratio).To(Equal(coretelemetry.DefaultTraceSampleRatio))
		interval,err:=coretelemetry.MetricExportInterval("20s")
		Expect(err).NotTo(HaveOccurred())
		Expect(interval).To(Equal(20*time.Second))
		Expect(coretelemetry.DefaultMetricExportInterval).To(Equal(15*time.Second))
	})
})
