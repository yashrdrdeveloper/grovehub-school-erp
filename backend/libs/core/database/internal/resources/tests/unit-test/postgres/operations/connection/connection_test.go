package connection_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres/operations/connection"
)

var _ = Describe("NewDB", func() {
	It("rejects malformed PostgreSQL connection strings without a live database", func() {
		_, err := connection.NewDB("postgres://%zz")

		Expect(err).To(HaveOccurred())
	})
})
