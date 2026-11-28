package config_test

import (
	"time"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/config"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/constants"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TraceSampleRatio", func() {
	It("parses valid ratios including both boundaries", func() {
		for input,expected:=range map[string]float64{"":1," ":1,"0":0,"0.25":0.25," 1 ":1}{
			ratio,err:=config.TraceSampleRatio(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(ratio).To(Equal(expected))
		}
	})
	It("returns defaults and errors for invalid values",func(){
		for _,input:=range []string{"-0.1","1.1","invalid","Inf"}{
			ratio,err:=config.TraceSampleRatio(input)
			Expect(err).To(MatchError("trace sample ratio must be between 0 and 1"))
			Expect(ratio).To(Equal(constants.DefaultTraceSampleRatio))
		}
	})
})

var _ = Describe("MetricExportInterval",func(){
	It("parses valid durations and defaults",func(){
		for input,expected:=range map[string]time.Duration{"":15*time.Second," 5s ":5*time.Second,"1m":time.Minute}{
			value,err:=config.MetricExportInterval(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal(expected))
		}
	})
	It("rejects invalid or too-short durations",func(){
		for _,input:=range []string{"4s","invalid","-1s"}{
			value,err:=config.MetricExportInterval(input)
			Expect(err).To(MatchError("metric export interval must be at least 5s"))
			Expect(value).To(Equal(constants.DefaultMetricExportInterval))
		}
	})
})
