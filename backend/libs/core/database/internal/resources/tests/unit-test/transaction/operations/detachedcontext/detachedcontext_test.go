package detachedcontext_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/database/internal/src/transaction/helpers"
	"grovehub/libs/core/database/internal/src/transaction/operations/detachedcontext"
)

var _ = Describe("DetachedContext", func() {
	It("retains values and removes cancellation and transaction state", func() {
		type key struct{}
		parent, cancel := context.WithCancel(context.WithValue(context.Background(), key{}, "trace"))
		child, _ := helpers.Bind(parent, nil)
		detached := detachedcontext.DetachedContext(child)
		cancel()

		Expect(detached.Value(key{})).To(Equal("trace"))
		Expect(detached.Err()).NotTo(HaveOccurred())
		_, ok := helpers.FromContext(detached)
		Expect(ok).To(BeFalse())
	})
})
