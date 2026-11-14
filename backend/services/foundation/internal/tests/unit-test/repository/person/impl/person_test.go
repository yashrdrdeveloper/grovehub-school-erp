package impl_test

import (
    "context"
    "errors"

    "github.com/DATA-DOG/go-sqlmock"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "gorm.io/gorm"

    database "grovehub/libs/core/tests/mocks/database"
    personimpl "grovehub/services/foundation/internal/app/repository/person/impl"
    repoerrors "grovehub/services/foundation/internal/app/repository/constants/errors"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    ids "grovehub/services/foundation/internal/tests/constants/person"
)

var db *gorm.DB
var sql sqlmock.Sqlmock

// @BeforeEach
func setupSQL() {
    var closeConnection func() error
    var err error
    db,sql,closeConnection,err=database.NewPostgres()
    Expect(err).NotTo(HaveOccurred())
    DeferCleanup(func(){Expect(sql.ExpectationsWereMet()).To(Succeed());Expect(closeConnection()).To(Succeed())})
}

// @Testing suite=FoundationPersonRepository
func PersonRepositorySpecs() {
    Describe("PersonRepository with isolated SQL expectations",func(){
        It("finds a person by GID without reaching PostgreSQL",func(){
            sql.ExpectQuery(`SELECT .*FROM .*persons.*`).WillReturnRows(sqlmock.NewRows([]string{"gid","first_name","last_name","status","version"}).AddRow(ids.GID,ids.FirstName,ids.LastName,"ACTIVE",1))
            p,err:=personimpl.NewPersonRepository(db).FindByGID(context.Background(),ids.GID)
            Expect(err).NotTo(HaveOccurred())
            Expect(p.GID).To(Equal(ids.GID))
            Expect(p.Status).To(Equal(fixtures.Entity().Status))
        })
        It("translates an absent person to the domain repository sentinel",func(){
            sql.ExpectQuery(`SELECT .*FROM .*persons.*`).WillReturnRows(sqlmock.NewRows([]string{"gid","first_name","last_name","status","version"}))
            _,err:=personimpl.NewPersonRepository(db).FindByGID(context.Background(),ids.GID)
            Expect(err).To(MatchError(repoerrors.ErrPersonNotFound))
        })
        It("preserves database query failures",func(){
            failure:=errors.New("simulated database read failure")
            sql.ExpectQuery(`SELECT .*FROM .*persons.*`).WillReturnError(failure)
            _,err:=personimpl.NewPersonRepository(db).FindByGID(context.Background(),ids.GID)
            Expect(err).To(MatchError(failure))
        })
        It("finds a person by idempotency key",func(){
            sql.ExpectQuery(`SELECT .*FROM .*persons.*`).WillReturnRows(sqlmock.NewRows([]string{"gid","first_name","last_name","status","version"}).AddRow(ids.GID,ids.FirstName,ids.LastName,"ACTIVE",1))
            p,err:=personimpl.NewPersonRepository(db).FindByIdempotencyKey(context.Background(),ids.IdempotencyKey)
            Expect(err).NotTo(HaveOccurred())
            Expect(p.GID).To(Equal(ids.GID))
        })
        It("maps an absent idempotency key to PersonNotFound",func(){
            sql.ExpectQuery(`SELECT .*FROM .*persons.*`).WillReturnRows(sqlmock.NewRows([]string{"gid"}))
            _,err:=personimpl.NewPersonRepository(db).FindByIdempotencyKey(context.Background(),ids.IdempotencyKey)
            Expect(err).To(MatchError(repoerrors.ErrPersonNotFound))
        })
        It("handles a failed idempotency query",func(){
            failure:=errors.New("database query failed")
            sql.ExpectQuery(`SELECT .*FROM .*persons.*`).WillReturnError(failure)
            _,err:=personimpl.NewPersonRepository(db).FindByIdempotencyKey(context.Background(),ids.IdempotencyKey)
            Expect(err).To(MatchError(failure))
        })
        It("rejects stale versions after an optimistic update",func(){
            sql.ExpectExec(`UPDATE .*persons.*`).WillReturnResult(sqlmock.NewResult(0,0))
            person:=fixtures.Entity()
            person.Version=2
            err:=personimpl.NewPersonRepository(db).Update(context.Background(),person)
            Expect(err).To(MatchError(repoerrors.ErrPersonVersionConflict))
        })
        It("updates a person with a current version",func(){
            sql.ExpectExec(`UPDATE .*persons.*`).WillReturnResult(sqlmock.NewResult(0,1))
            person:=fixtures.Entity()
            person.Version=2
            Expect(personimpl.NewPersonRepository(db).Update(context.Background(),person)).To(Succeed())
        })
        It("propagates a write failure",func(){
            failure:=errors.New("simulated update error")
            sql.ExpectExec(`UPDATE .*persons.*`).WillReturnError(failure)
            person:=fixtures.Entity()
            person.Version=2
            Expect(personimpl.NewPersonRepository(db).Update(context.Background(),person)).To(MatchError(failure))
        })
    })
}
