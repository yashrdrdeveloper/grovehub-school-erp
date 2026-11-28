package service_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corepipeline "grovehub/libs/core/pipeline"
)

type fakeUnit struct {
	completed bool
	err       error
}

func (u *fakeUnit) Run(ctx context.Context, fn func(context.Context) error) error {
	err := fn(ctx)
	u.completed = true
	if u.err != nil {
		return u.err
	}
	return err
}

var _ = Describe("Transactional", func() {
	It("runs recovery only after a failed transaction has returned", func() {
		sentinel := errors.New("duplicate")
		unit := &fakeUnit{}
		value, created, err := corepipeline.RunService(context.Background(), new(int), func(ctx context.Context, s *int) corepipeline.Outcome[string] {
			handler := corepipeline.Transactional[int, string](unit, func(context.Context, *int) corepipeline.Outcome[string] { return corepipeline.Fail[string](sentinel) })
			return corepipeline.Recover(ctx, s, handler(ctx, s), corepipeline.OnError[int, string](sentinel, func(context.Context, *int, error) corepipeline.Outcome[string] {
				Expect(unit.completed).To(BeTrue())
				return corepipeline.Return("existing", false)
			}))
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("existing"))
		Expect(created).To(BeFalse())
	})

	It("requires a unit of work and preserves successful outcomes", func() {
		handler := func(context.Context, *int) corepipeline.Outcome[string] { return corepipeline.Return("new", true) }
		_, _, err := corepipeline.RunService(context.Background(), new(int), corepipeline.Transactional[int, string](nil, handler))
		Expect(err).To(MatchError(corepipeline.ErrMissingUnitOfWork))
		value, created, err := corepipeline.RunService(context.Background(), new(int), corepipeline.Transactional[int, string](&fakeUnit{}, handler))
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("new"))
		Expect(created).To(BeTrue())
	})

	It("keeps unmatched failures and permits fallback recovery", func() {
		sentinel := errors.New("unexpected")
		_, _, err := corepipeline.RunService(context.Background(), new(int), func(ctx context.Context, s *int) corepipeline.Outcome[string] {
			return corepipeline.Recover(ctx, s, corepipeline.Fail[string](sentinel), corepipeline.Recovery[int, string]{})
		})
		Expect(err).To(MatchError(sentinel))
		value, _, err := corepipeline.RunService(context.Background(), new(int), func(ctx context.Context, s *int) corepipeline.Outcome[string] {
			return corepipeline.Recover(ctx, s, corepipeline.Fail[string](sentinel), corepipeline.OtherwiseError[int, string](func(context.Context, *int, error) corepipeline.Outcome[string] {
				return corepipeline.Return("fallback", false)
			}))
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("fallback"))
	})
})

var _ = Describe("Sequence", func() {
	It("stops on error and returns a result only after success", func() {
		sentinel := errors.New("write failed")
		_, _, err := corepipeline.RunService(context.Background(), new(int), func(ctx context.Context, _ *int) corepipeline.Outcome[string] {
			return corepipeline.Sequence(ctx, "done", true, func(context.Context) error { return sentinel }, func(context.Context) error { Fail("unexpected write"); return nil })
		})
		Expect(err).To(MatchError(sentinel))
		value, created, err := corepipeline.RunService(context.Background(), new(int), func(ctx context.Context, s *int) corepipeline.Outcome[string] {
			return corepipeline.Recover(ctx, s, corepipeline.Sequence(ctx, "done", true, func(context.Context) error { return nil }))
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("done"))
		Expect(created).To(BeTrue())
	})
})
