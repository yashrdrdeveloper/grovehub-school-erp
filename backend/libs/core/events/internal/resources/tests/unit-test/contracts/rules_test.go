package contracts_test

import (
	"context"
	"grovehub/libs/core/events/internal/src/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rules", func() {
	It("keeps the consumer-defined message type and the stable incomplete-rules error", func() {
		rules := contracts.Rules[int, string]{
			ShouldEmit: func(value int) bool { return value > 0 },
			Build: func(_ context.Context, _ int) (string, error) { return "message", nil },
			Enqueue: func(_ context.Context, _ string) error { return nil },
		}
		Expect(rules.ShouldEmit(1)).To(BeTrue())
		message, err := rules.Build(context.Background(), 1)
		Expect(err).NotTo(HaveOccurred())
		Expect(message).To(Equal("message"))
		Expect(rules.Enqueue(context.Background(), message)).To(Succeed())
		Expect(contracts.ErrIncompleteRules).To(MatchError("event pipeline requires ShouldEmit, Build and Enqueue"))
	})
})
