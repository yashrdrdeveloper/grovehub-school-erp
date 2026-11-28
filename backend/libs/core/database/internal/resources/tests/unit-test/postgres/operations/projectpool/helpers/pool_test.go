package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres/operations/projectpool/constants"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres/operations/projectpool/helpers"
)

var _ = Describe("ConfigurePool", func() {
	It("sets the maximum number of open connections", func() {
		db, mock, err := sqlmock.New()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			mock.ExpectClose()
			Expect(db.Close()).To(Succeed())
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		})

		helpers.ConfigurePool(db)

		Expect(db.Stats().MaxOpenConnections).To(Equal(constants.MaxOpenConns))
	})
})
