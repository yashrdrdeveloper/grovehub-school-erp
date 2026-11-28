package events_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreevents "grovehub/libs/core/events"
)

var _ = Describe("Run", func() {
	It("preserves public rules, eligibility, delivery and error identity", func() {
		type message struct{ EventID string }
		calls := 0
		sentinel := errors.New("outbox unavailable")
		rules := coreevents.Rules[int, message]{
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
		Expect(coreevents.Run(context.Background(), 0, rules)).To(Succeed())
		Expect(calls).To(BeZero())
		Expect(coreevents.Run(context.Background(), 1, rules)).To(BeIdenticalTo(sentinel))
		Expect(calls).To(Equal(2))
		Expect(coreevents.Run(context.Background(), 1, coreevents.Rules[int, message]{})).To(BeIdenticalTo(coreevents.ErrIncompleteRules))
	})
})
