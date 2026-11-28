package postgres_test

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
)

var _ = Describe("OpenProjectPool", func() {
	It("permits each project role only on its own database", func() {
		school := os.Getenv("PROJECT_SCHOOL_TEST_DSN")
		internal := os.Getenv("PROJECT_INTERNAL_TEST_DSN")
		denied := os.Getenv("PROJECT_CROSS_TEST_DSN")

		if school == "" && internal == "" && denied == "" {
			Skip("requires disposable project databases")
		}

		Expect(school).NotTo(BeEmpty())
		Expect(internal).NotTo(BeEmpty())
		Expect(denied).NotTo(BeEmpty())

		for _, dsn := range []string{school, internal} {
			db, err := coredatabase.OpenProjectPool(context.Background(), dsn)
			Expect(err).NotTo(HaveOccurred())
			DeferCleanup(db.Close)
			Expect(db.PingContext(context.Background())).To(Succeed())
		}

		db, err := coredatabase.OpenProjectPool(context.Background(), denied)
		if db != nil {
			_ = db.Close()
		}

		Expect(err).To(HaveOccurred())
		Expect(db).To(BeNil())
	})
})
