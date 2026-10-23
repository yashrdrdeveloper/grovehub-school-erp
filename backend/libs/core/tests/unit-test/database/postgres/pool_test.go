package postgres_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corepostgres "grovehub/libs/core/database/postgres"
)

var _ = Describe("PostgreSQL pool", func() {
	It("creates a pool from a valid database URL without opening a connection", func() {
		pool, err := corepostgres.NewPool(
			context.Background(),
			"postgres://grovehub:grovehub@localhost:5432/grovehub?sslmode=disable",
		)

		Expect(err).NotTo(HaveOccurred())
		Expect(pool).NotTo(BeNil())

		pool.Close()
	})

	It("returns an error for an invalid database URL", func() {
		pool, err := corepostgres.NewPool(
			context.Background(),
			"postgres://grovehub:grovehub@localhost:not-a-port/grovehub",
		)

		Expect(err).To(HaveOccurred())
		Expect(pool).To(BeNil())
	})
})
