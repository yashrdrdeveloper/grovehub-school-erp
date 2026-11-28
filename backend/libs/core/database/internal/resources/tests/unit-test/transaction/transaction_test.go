package transaction_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/database/internal/src/transaction"
)

var _ = Describe("Run", func() {
	It("exposes the same missing-database error as the transaction operation", func() {
		Expect(transaction.Run(context.Background(), nil, nil)).To(MatchError(transaction.ErrMissingDatabase))
	})
})

var _ = Describe("DetachedContext", func() {
	It("exposes a context without its original cancellation", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		Expect(transaction.DetachedContext(ctx).Err()).NotTo(HaveOccurred())
	})
})
