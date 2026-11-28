package pipeline_test

import (
	"context"
	"errors"
	"grovehub/libs/core/events/internal/src/contracts"
	"grovehub/libs/core/events/internal/src/operations/pipeline"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Run", func() {
	It("rejects every missing callback before invoking any operation", func() {
		valid := contracts.Rules[int, string]{
			ShouldEmit: func(int) bool { Fail("unexpected eligibility check"); return false },
			Build:      func(context.Context, int) (string, error) { Fail("unexpected build"); return "", nil },
			Enqueue:    func(context.Context, string) error { Fail("unexpected enqueue"); return nil },
		}
		for _, rules := range []contracts.Rules[int, string]{
			{},
			{Build: valid.Build, Enqueue: valid.Enqueue},
			{ShouldEmit: valid.ShouldEmit, Enqueue: valid.Enqueue},
			{ShouldEmit: valid.ShouldEmit, Build: valid.Build},
		} {
			Expect(errors.Is(pipeline.Run(context.Background(), 1, rules), contracts.ErrIncompleteRules)).To(BeTrue())
		}
	})

	It("does not build or enqueue an ineligible event", func() {
		rules := contracts.Rules[int, string]{
			ShouldEmit: func(value int) bool { return value > 0 },
			Build:      func(context.Context, int) (string, error) { Fail("unexpected build"); return "", nil },
			Enqueue:    func(context.Context, string) error { Fail("unexpected enqueue"); return nil },
		}
		Expect(pipeline.Run(context.Background(), 0, rules)).To(Succeed())
	})

	It("passes the original context, builds once and enqueues once in order", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		calls := []string{}
		rules := contracts.Rules[int, string]{
			ShouldEmit: func(value int) bool { calls = append(calls, "eligible"); return value == 7 },
			Build: func(received context.Context, value int) (string, error) {
				Expect(received).To(BeIdenticalTo(ctx))
				Expect(value).To(Equal(7))
				calls = append(calls, "build")
				return "message", nil
			},
			Enqueue: func(received context.Context, message string) error {
				Expect(received).To(BeIdenticalTo(ctx))
				Expect(message).To(Equal("message"))
				calls = append(calls, "enqueue")
				return nil
			},
		}
		Expect(pipeline.Run(ctx, 7, rules)).To(Succeed())
		Expect(calls).To(Equal([]string{"eligible", "build", "enqueue"}))
	})

	It("returns the original build failure and skips enqueue", func() {
		sentinel := errors.New("build failed")
		rules := contracts.Rules[int, string]{
			ShouldEmit: func(int) bool { return true },
			Build:      func(context.Context, int) (string, error) { return "", sentinel },
			Enqueue:    func(context.Context, string) error { Fail("unexpected enqueue"); return nil },
		}
		Expect(pipeline.Run(context.Background(), 1, rules)).To(BeIdenticalTo(sentinel))
	})

	It("returns the original enqueue failure", func() {
		sentinel := errors.New("enqueue failed")
		rules := contracts.Rules[int, string]{
			ShouldEmit: func(int) bool { return true },
			Build:      func(context.Context, int) (string, error) { return "message", nil },
			Enqueue:    func(context.Context, string) error { return sentinel },
		}
		Expect(pipeline.Run(context.Background(), 1, rules)).To(BeIdenticalTo(sentinel))
	})
})
