package validators_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres/operations/projectpool/validators"
)

var _ = Describe("ProjectConfig", func() {
	DescribeTable("rejects invalid or overridden project identities",
		func(dsn string) {
			config, err := validators.ProjectConfig(dsn)

			Expect(err).To(HaveOccurred())
			Expect(config).To(BeNil())
		},
		Entry("empty", ""),
		Entry("keyword DSN", "host=localhost dbname=school"),
		Entry("missing user", "postgres://localhost/school"),
		Entry("missing database", "postgres://school@localhost"),
		Entry("database override", "postgres://school@localhost/school?dbname=platform"),
		Entry("user override", "postgres://school@localhost/school?user=platform"),
		Entry("service override", "postgres://school@localhost/school?service=platform"),
		Entry("host override", "postgres://school@localhost/school?host=another"),
		Entry("fragment", "postgres://school@localhost/school#fragment"),
		Entry("invalid escape", "postgres://school@localhost/school?sslmode=%"),
	)

	It("accepts explicit database and role", func() {
		config, err := validators.ProjectConfig("postgres://school:password@127.0.0.1:5432/school?sslmode=disable")

		Expect(err).NotTo(HaveOccurred())
		Expect(config.Database).To(Equal("school"))
		Expect(config.User).To(Equal("school"))
	})
})
