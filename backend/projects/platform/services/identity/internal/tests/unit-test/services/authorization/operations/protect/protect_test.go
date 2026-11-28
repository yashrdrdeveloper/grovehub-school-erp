package protect_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/authorization/operations/protect"
)

func TestAuthorityGuard(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Authority Guard Suite")
}

type locker struct {
	tenant string
	calls int
}

func (lock *locker) Run(ctx context.Context, tenantID string, operation func(context.Context) error) error {
	lock.tenant = tenantID
	lock.calls++
	return operation(ctx)
}

type reader struct {
	counts []int64
	err error
	calls int
}

func (read *reader) Count(ctx context.Context, tenantID string) (int64, error) {
	read.calls++
	if read.err != nil {
		return 0, read.err
	}
	return read.counts[read.calls-1], nil
}

var _ = Describe("Run", func() {
	It("allows an authority mutation only while an administrator remains", func() {
		lock := &locker{}
		read := &reader{counts: []int64{2, 1}}
		calls := 0
		err := protect.New(lock, read).Run(context.Background(), "school", func(ctx context.Context) error {
			calls++
			return nil
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(lock.tenant).To(Equal("school"))
		Expect(calls).To(Equal(1))
		Expect(read.calls).To(Equal(2))
	})
	It("rejects the last removal and refuses an unbootstrapped tenant", func() {
		calls := 0
		operation := func(ctx context.Context) error { calls++; return nil }
		err := protect.New(&locker{}, &reader{counts: []int64{1, 0}}).Run(context.Background(), "school", operation)
		Expect(err).To(MatchError(authority.ErrLastAdministrator))
		err = protect.New(&locker{}, &reader{counts: []int64{0}}).Run(context.Background(), "school", operation)
		Expect(err).To(MatchError(authority.ErrBootstrapRequired))
		Expect(calls).To(Equal(1))
	})
	It("fails closed for dependency, scope and operation errors", func() {
		failure := errors.New("failed")
		operation := func(ctx context.Context) error { return failure }
		Expect(protect.New(nil, nil).Run(context.Background(), "school", operation)).To(MatchError(authorization.ErrUnavailable))
		Expect(protect.New(&locker{}, &reader{}).Run(context.Background(), "", operation)).To(MatchError(authority.ErrScopeRequired))
		Expect(protect.New(&locker{}, &reader{err: failure}).Run(context.Background(), "school", operation)).To(MatchError(failure))
		read := &reader{counts: []int64{1}}
		Expect(protect.New(&locker{}, read).Run(context.Background(), "school", operation)).To(MatchError(failure))
		Expect(read.calls).To(Equal(1))
	})
})
