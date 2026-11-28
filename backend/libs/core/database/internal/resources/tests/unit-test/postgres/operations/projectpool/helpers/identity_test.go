package helpers_test

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/database/internal/src/postgres/operations/projectpool/helpers"
)

var _ = Describe("MatchesIdentity", func() {
	var (
		db   *sql.DB
		mock sqlmock.Sqlmock
	)

	BeforeEach(func() {
		var err error
		db, mock, err = sqlmock.New()
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func() {
			Expect(mock.ExpectationsWereMet()).To(Succeed())
			mock.ExpectClose()
			Expect(db.Close()).To(Succeed())
		})
	})

	DescribeTable("checks both database and role",
		func(database, role string, expected bool) {
			mock.ExpectQuery("SELECT current_database\\(\\), current_user").
				WillReturnRows(sqlmock.NewRows([]string{"database", "role"}).AddRow(database, role))

			Expect(helpers.MatchesIdentity(context.Background(), db, "school", "school_role")).To(Equal(expected))
		},
		Entry("matching", "school", "school_role", true),
		Entry("wrong database", "other", "school_role", false),
		Entry("wrong role", "school", "other", false),
	)

	It("fails closed on query errors", func() {
		mock.ExpectQuery("SELECT current_database\\(\\), current_user").WillReturnError(errors.New("connection failed"))

		Expect(helpers.MatchesIdentity(context.Background(), db, "school", "school_role")).To(BeFalse())
	})

})
