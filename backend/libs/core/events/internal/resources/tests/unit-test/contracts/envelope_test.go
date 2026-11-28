package contracts_test

import (
	"encoding/json"
	"time"
	"grovehub/libs/core/events/internal/src/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Envelope", func() {
	It("preserves wire names, generic data and optional IDs", func() {
		value := contracts.Envelope[string]{
			EventID: "event", EventType: "changed", Version: 2, AggregateID: "aggregate",
			OccurredAt: time.Unix(0, 0).UTC(), Data: "payload",
		}
		raw, err := json.Marshal(value)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(raw)).To(MatchJSON(`{"eventId":"event","eventType":"changed","version":2,"aggregateId":"aggregate","occurredAt":"1970-01-01T00:00:00Z","data":"payload"}`))
		var decoded contracts.Envelope[string]
		Expect(json.Unmarshal(raw, &decoded)).To(Succeed())
		Expect(decoded).To(Equal(value))
	})
})
