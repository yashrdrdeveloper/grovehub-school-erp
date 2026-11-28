package projectpool_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres/operations/projectpool"
)

var _ = Describe("OpenProjectPool", func() {
	It("rejects an identity override without opening a connection or exposing credentials", func() {
		pool, err := projectpool.OpenProjectPool(context.Background(), "postgres://school:private-password@localhost/school?user=other")

		Expect(pool).To(BeNil())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).NotTo(ContainSubstring("private-password"))
	})
})
