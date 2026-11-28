package causation_test

import (
	"context"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/causation"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/internal/src/operations/correlation"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("WithID", func() {
	It("propagates causation independently of correlation and the parent context", func() {
		parent := correlation.WithID(context.Background(), "request")
		ctx := causation.WithID(parent, "command")
		Expect(causation.Resolve(ctx, "")).To(Equal("command"))
		Expect(correlation.Resolve(ctx, "")).To(Equal("request"))
		Expect(causation.Resolve(parent, "")).To(BeEmpty())
	})
})

var _ = Describe("Resolve", func() {
	It("prefers metadata and preserves unknown IDs", func() {
		ctx := causation.WithID(context.Background(), "context")
		Expect(causation.Resolve(ctx, "metadata")).To(Equal("metadata"))
		Expect(causation.Resolve(causation.WithID(ctx, ""), "")).To(BeEmpty())
		Expect(causation.Resolve(context.Background(), "")).To(BeEmpty())
	})
})
