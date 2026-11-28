package service_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corepipeline "grovehub/libs/core/pipeline"
)

var _ = Describe("Decide", func() {
	It("chooses only the first matching rule and supports a fallback", func() {
		for _, state := range []int{1, 2} {
			value, _, err := corepipeline.RunService(context.Background(), &state, func(ctx context.Context, s *int) corepipeline.Outcome[string] {
				return corepipeline.Decide(ctx, s, corepipeline.Rule[int, string]{},
					corepipeline.When(func(s *int) bool { return *s == 1 }, func(context.Context, *int) corepipeline.Outcome[string] { return corepipeline.Return("first", true) }),
					corepipeline.Otherwise(func(context.Context, *int) corepipeline.Outcome[string] {
						return corepipeline.Return("fallback", false)
					}),
				)
			})
			Expect(err).NotTo(HaveOccurred())
			if state == 1 {
				Expect(value).To(Equal("first"))
			} else {
				Expect(value).To(Equal("fallback"))
			}
		}
		_, _, err := corepipeline.RunService(context.Background(), new(int), func(ctx context.Context, s *int) corepipeline.Outcome[string] {
			return corepipeline.Decide[int, string](ctx, s)
		})
		Expect(err).To(MatchError(corepipeline.ErrNoMatchingRule))
	})
})

var _ = Describe("Lookup", func() {
	It("distinguishes success, missing and unexpected errors", func() {
		missing := errors.New("missing")
		unexpected := errors.New("storage error")
		for _, item := range []struct {
			err  error
			want string
		}{{nil, "found"}, {missing, "not found"}, {unexpected, "failed"}} {
			value, _, err := corepipeline.RunService(context.Background(), new(int), func(context.Context, *int) corepipeline.Outcome[string] {
				return corepipeline.Lookup(1, item.err, missing, func(int) corepipeline.Outcome[string] { return corepipeline.Return("found", false) }, func() corepipeline.Outcome[string] { return corepipeline.Return("not found", false) }, func(error) corepipeline.Outcome[string] { return corepipeline.Return("failed", false) })
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal(item.want))
		}
	})
})
