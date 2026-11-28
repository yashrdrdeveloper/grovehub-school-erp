package contracts_test

import (
	"time"
	"grovehub/libs/core/events/internal/src/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metadata", func() {
	It("retains supplied event identity, time and schema version", func() {
		at := time.Unix(10, 0).UTC()
		value := contracts.Metadata{
			EventID: "event", EventType: "changed", Version: 3, AggregateID: "aggregate",
			OccurredAt: at, CorrelationID: "request", CausationID: "command",
		}
		Expect(value.EventID).To(Equal("event"))
		Expect(value.EventType).To(Equal("changed"))
		Expect(value.Version).To(Equal(3))
		Expect(value.AggregateID).To(Equal("aggregate"))
		Expect(value.OccurredAt).To(Equal(at))
		Expect(value.CorrelationID).To(Equal("request"))
		Expect(value.CausationID).To(Equal("command"))
	})
})
