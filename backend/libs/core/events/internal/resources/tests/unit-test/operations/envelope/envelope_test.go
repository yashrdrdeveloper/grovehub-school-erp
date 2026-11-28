package envelope_test

import (
	"context"
	"time"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/envelope"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/causation"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/correlation"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("New", func() {
	It("copies metadata and payload and resolves context IDs", func() {
		at := time.Unix(20, 0).UTC()
		ctx := causation.WithID(correlation.WithID(context.Background(), "request"), "command")
		metadata := contracts.Metadata{EventID: "event", EventType: "changed", Version: 2, AggregateID: "aggregate", OccurredAt: at}
		payload := []int{1, 2}
		Expect(envelope.New(ctx, metadata, payload)).To(Equal(contracts.Envelope[[]int]{
			EventID: "event", EventType: "changed", Version: 2, AggregateID: "aggregate", OccurredAt: at,
			CorrelationID: "request", CausationID: "command", Data: payload,
		}))
		metadata.CorrelationID = "explicit-request"
		metadata.CausationID = "explicit-command"
		event := envelope.New(ctx, metadata, payload)
		Expect(event.CorrelationID).To(Equal("explicit-request"))
		Expect(event.CausationID).To(Equal("explicit-command"))
	})

	It("preserves zero values and nil payloads without inventing metadata", func() {
		Expect(envelope.New(context.Background(), contracts.Metadata{}, (*int)(nil))).To(Equal(contracts.Envelope[*int]{}))
	})
})
