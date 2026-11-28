package postgres_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/database/internal/src/postgres"
)

var _ = Describe("OpenProjectPool", func() {
	It("forwards URI validation failures through the private facade", func() {
		pool, err := postgres.OpenProjectPool(context.Background(), "")

		Expect(pool).To(BeNil())
		Expect(err).To(MatchError(ContainSubstring("explicit project PostgreSQL URI")))
	})
})

var _ = Describe("NewDB", func() {
	It("returns an error for a malformed connection URI", func() {
		_, err := postgres.NewDB("postgres://%zz")

		Expect(err).To(HaveOccurred())
	})
})
