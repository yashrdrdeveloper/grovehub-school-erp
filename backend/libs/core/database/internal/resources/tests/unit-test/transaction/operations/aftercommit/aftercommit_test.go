package aftercommit_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/database/internal/src/transaction/helpers"
	"grovehub/libs/core/database/internal/src/transaction/operations/aftercommit"
)

var _ = Describe("AfterCommit", func() {
	It("registers a callback only on a transaction and defers execution", func() {
		ctx := context.Background()
		child, state := helpers.Bind(ctx, nil)
		called := false
		callback := func() { called = true }

		Expect(aftercommit.AfterCommit(ctx, callback)).To(BeFalse())
		Expect(aftercommit.AfterCommit(child, callback)).To(BeTrue())
		Expect(called).To(BeFalse())

		state.PublishCallbacks()

		Expect(called).To(BeTrue())
	})
})
