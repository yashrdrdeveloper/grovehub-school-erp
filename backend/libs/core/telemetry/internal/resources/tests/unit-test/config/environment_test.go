package config_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"time"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment configuration",func(){
	It("uses trimmed values or explicit defaults",func(){
		GinkgoT().Setenv("TELEMETRY_TEST_VALUE","   ")
		Expect(config.EnvOrDefault("TELEMETRY_TEST_VALUE","fallback")).To(Equal("fallback"))
		GinkgoT().Setenv("TELEMETRY_TEST_VALUE"," value ")
		Expect(config.EnvOrDefault("TELEMETRY_TEST_VALUE","fallback")).To(Equal("value"))
		Expect(config.DefaultInstanceID()).To(HaveSuffix(fmt.Sprintf("-%d",os.Getpid())))
	})
	It("resolves signal-specific and common export endpoints",func(){
		for _,key:=range []string{"OTEL_EXPORTER_OTLP_ENDPOINT","OTEL_EXPORTER_OTLP_TRACES_ENDPOINT","OTEL_EXPORTER_OTLP_METRICS_ENDPOINT"}{GinkgoT().Setenv(key,"")}
		Expect(config.TraceExportConfigured()).To(BeFalse())
		Expect(config.MetricExportConfigured()).To(BeFalse())
		GinkgoT().Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT","http://localhost:4318")
		Expect(config.TraceExportConfigured()).To(BeTrue())
		Expect(config.MetricExportConfigured()).To(BeFalse())
		GinkgoT().Setenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT","http://localhost:4318")
		Expect(config.MetricExportConfigured()).To(BeTrue())
		GinkgoT().Setenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT","")
		GinkgoT().Setenv("OTEL_EXPORTER_OTLP_METRICS_ENDPOINT","")
		GinkgoT().Setenv("OTEL_EXPORTER_OTLP_ENDPOINT","http://localhost:4318")
		Expect(config.TraceExportConfigured()).To(BeTrue())
		Expect(config.MetricExportConfigured()).To(BeTrue())
	})
	It("falls back to defaults and warns for invalid settings",func(){
		var output bytes.Buffer
		logger:=slog.New(slog.NewTextHandler(&output,nil))
		GinkgoT().Setenv("GROVEHUB_TRACE_SAMPLE_RATIO","0.5")
		GinkgoT().Setenv("GROVEHUB_METRIC_EXPORT_INTERVAL","20s")
		Expect(config.ConfiguredTraceSampleRatio(logger)).To(Equal(0.5))
		Expect(config.ConfiguredMetricExportInterval(logger)).To(Equal(20*time.Second))
		GinkgoT().Setenv("GROVEHUB_TRACE_SAMPLE_RATIO","invalid")
		GinkgoT().Setenv("GROVEHUB_METRIC_EXPORT_INTERVAL","1s")
		Expect(config.ConfiguredTraceSampleRatio(logger)).To(Equal(1.0))
		Expect(config.ConfiguredMetricExportInterval(logger)).To(Equal(15*time.Second))
		Expect(output.String()).To(ContainSubstring("invalid GROVEHUB_TRACE_SAMPLE_RATIO"))
		Expect(output.String()).To(ContainSubstring("invalid GROVEHUB_METRIC_EXPORT_INTERVAL"))
	})
})
