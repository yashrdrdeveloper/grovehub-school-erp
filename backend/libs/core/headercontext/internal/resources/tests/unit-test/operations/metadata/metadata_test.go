package metadata_test

import (
	"context"
	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext/internal/src/operations/metadata"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithValues", func() {
	It("copies trusted metadata, preserves parent state and propagates event correlation", func() {
		parent := context.Background()
		values := contracts.Values{IdempotencyKey: "key", CorrelationID: "corr", RequestID: "req"}
		ctx := metadata.WithValues(parent, values)
		values.IdempotencyKey = "changed"
		Expect(metadata.FromContext(ctx)).To(Equal(contracts.Values{IdempotencyKey: "key", CorrelationID: "corr", RequestID: "req"}))
		Expect(metadata.FromContext(parent)).To(Equal(contracts.Values{}))
		Expect(coreevents.NewEnvelope(ctx, coreevents.Metadata{}, "data").CorrelationID).To(Equal("corr"))
	})

	It("does not replace an existing event correlation when the new value is empty", func() {
		parent := coreevents.WithCorrelationID(context.Background(), "existing")
		ctx := metadata.WithValues(parent, contracts.Values{})
		Expect(coreevents.NewEnvelope(ctx, coreevents.Metadata{}, "data").CorrelationID).To(Equal("existing"))
	})
})

var _ = Describe("WithoutIdempotencyKey", func() {
	It("removes only the inherited key without modifying the parent", func() {
		ctx := metadata.WithValues(context.Background(), contracts.Values{IdempotencyKey: "key", CorrelationID: "corr", RequestID: "req"})
		child := metadata.WithoutIdempotencyKey(ctx)
		Expect(metadata.FromContext(child)).To(Equal(contracts.Values{CorrelationID: "corr", RequestID: "req"}))
		Expect(metadata.FromContext(ctx).IdempotencyKey).To(Equal("key"))
	})
})
