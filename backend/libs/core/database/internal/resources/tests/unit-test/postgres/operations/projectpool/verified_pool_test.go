package projectpool_test

import (
	"context"
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/postgres/operations/projectpool"
)

var _=Describe("Verified project pool",func(){
	const dsn="postgres://school:secret@127.0.0.1:5432/school?sslmode=disable"
	It("opens a pool only when the server confirms the expected database and role",func(){
		db,mock,err:=sqlmock.New()
		Expect(err).NotTo(HaveOccurred())
		mock.ExpectQuery("SELECT current_database").WillReturnRows(sqlmock.NewRows([]string{"database","role"}).AddRow("school","school"))
		pool,err:=projectpool.OpenVerifiedPool(context.Background(),dsn,func(config pgx.ConnConfig)*sql.DB{
			Expect(config.Database).To(Equal("school"))
			Expect(config.User).To(Equal("school"))
			return db
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(pool).To(BeIdenticalTo(db))
		mock.ExpectClose()
		Expect(pool.Close()).To(Succeed())
		Expect(mock.ExpectationsWereMet()).To(Succeed())
	})
	It("closes mismatched pools without revealing the supplied credentials",func(){
		db,mock,err:=sqlmock.New()
		Expect(err).NotTo(HaveOccurred())
		mock.ExpectQuery("SELECT current_database").WillReturnRows(sqlmock.NewRows([]string{"database","role"}).AddRow("school","other"))
		mock.ExpectClose()
		pool,err:=projectpool.OpenVerifiedPool(context.Background(),dsn,func(pgx.ConnConfig)*sql.DB{return db})
		Expect(pool).To(BeNil())
		Expect(err).To(MatchError("project PostgreSQL connection or identity verification failed"))
		Expect(mock.ExpectationsWereMet()).To(Succeed())
	})
	It("uses the real driver but does not connect when the context is cancelled",func(){
		ctx,cancel:=context.WithCancel(context.Background())
		cancel()
		pool,err:=projectpool.OpenProjectPool(ctx,dsn)
		Expect(pool).To(BeNil())
		Expect(err).To(MatchError("project PostgreSQL connection or identity verification failed"))
	})
})
