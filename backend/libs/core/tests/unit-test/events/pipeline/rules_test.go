package pipeline_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events/pipeline"
)

func TestPipeline(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Public Event Pipeline Suite")
}

var _ = Describe("Run", func() {
	It("preserves public rules, eligibility, delivery and error identity", func() {
		type message struct{ EventID string }
		calls := 0
		sentinel := errors.New("outbox unavailable")
		rules := pipeline.Rules[int, message]{
			ShouldEmit: func(value int) bool { return value > 0 },
			Build: func(context.Context, int) (message, error) {
				calls++
				return message{EventID: "event"}, nil
			},
			Enqueue: func(_ context.Context, value message) error {
				Expect(value.EventID).To(Equal("event"))
				calls++
				return sentinel
			},
		}
		Expect(pipeline.Run(context.Background(), 0, rules)).To(Succeed())
		Expect(calls).To(BeZero())
		Expect(pipeline.Run(context.Background(), 1, rules)).To(BeIdenticalTo(sentinel))
		Expect(calls).To(Equal(2))
		Expect(pipeline.Run(context.Background(), 1, pipeline.Rules[int, message]{})).To(BeIdenticalTo(pipeline.ErrIncompleteRules))
	})
})
