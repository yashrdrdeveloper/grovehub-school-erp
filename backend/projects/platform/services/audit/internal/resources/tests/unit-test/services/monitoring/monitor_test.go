package monitoring_test

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/services/audit/internal/app/services/monitoring"
)

func TestMonitoring(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Audit monitoring suite")
}

var _ = Describe("Ready", func() {
	It("requires a fresh dependency check and every subscription", func() {
		monitor, err := monitoring.New(map[string]string{"one": "event.one", "two": "event.two"})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(monitor.Close)
		Expect(monitor.Ready()).To(BeFalse())
		monitor.Probe(time.Now(), true)
		monitor.Subscription("one", true)
		Expect(monitor.Ready()).To(BeFalse())
		monitor.Subscription("two", true)
		Expect(monitor.Ready()).To(BeTrue())
		monitor.Delivery(context.Background(), "one", "accepted", time.Millisecond, errors.New("failure"))
		Expect(monitor.Ready()).To(BeFalse())
		monitor.Delivery(context.Background(), "one", "invalid_event", time.Millisecond, nil)
		Expect(monitor.Ready()).To(BeTrue())
		monitor.Interrupted(context.Background(), "two")
		Expect(monitor.Ready()).To(BeFalse())
		monitor.Subscription("two", true)
		monitor.Probe(time.Now().Add(-time.Minute), true)
		Expect(monitor.Ready()).To(BeFalse())
		monitor.Probe(time.Now(), false)
		Expect(monitor.Ready()).To(BeFalse())
		monitor.Probe(time.Now(), true)
		Expect(monitor.Ready()).To(BeTrue())
		monitor.Close()
		Expect(monitor.Ready()).To(BeFalse())
	})
})
