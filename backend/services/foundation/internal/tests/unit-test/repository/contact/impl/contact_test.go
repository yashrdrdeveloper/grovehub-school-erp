package impl_test

import (
    "context"
    "errors"

    "github.com/DATA-DOG/go-sqlmock"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "gorm.io/gorm"

    database "grovehub/libs/core/tests/mocks/database"
    contactimpl "grovehub/services/foundation/internal/app/repository/contact/impl"
    "grovehub/services/foundation/internal/app/entity"
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
    DeferCleanup(func(){Expect(sql.ExpectationsWereMet()).To(Succeed());sql.ExpectClose();Expect(closeConnection()).To(Succeed())})
}

// @Testing suite=FoundationContactRepository
func ContactRepositorySpecs() {
    Describe("ContactRepository with isolated SQL expectations",func(){
        It("does not query for an empty set of person IDs",func(){
            values,err:=contactimpl.NewContactRepository(db).FindByPersonGIDs(context.Background(),nil)
            Expect(err).NotTo(HaveOccurred())
            Expect(values).To(BeEmpty())
        })
        It("reads and groups contact rows by person",func(){
            rows:=sqlmock.NewRows([]string{"person_gid","type","value","normalized_value","is_primary","is_verified"}).AddRow(ids.GID,"EMAIL",ids.Email,ids.Email,true,false)
            sql.ExpectQuery(`SELECT .*FROM .*person_contacts.*`).WillReturnRows(rows)
            values,err:=contactimpl.NewContactRepository(db).FindByPersonGIDs(context.Background(),[]string{ids.GID})
            Expect(err).NotTo(HaveOccurred())
            Expect(values).To(HaveKey(ids.GID))
            Expect(values[ids.GID]).To(HaveLen(1))
            Expect(values[ids.GID][0].NormalizedValue).To(Equal(ids.Email))
        })
        It("returns a single person's contacts through FindByPersonGID",func(){
            rows:=sqlmock.NewRows([]string{"person_gid","type","value","normalized_value"}).AddRow(ids.GID,"EMAIL",ids.Email,ids.Email)
            sql.ExpectQuery(`SELECT .*FROM .*person_contacts.*`).WillReturnRows(rows)
            values,err:=contactimpl.NewContactRepository(db).FindByPersonGID(context.Background(),ids.GID)
            Expect(err).NotTo(HaveOccurred())
            Expect(values).To(HaveLen(1))
        })
        It("preserves a contact query error",func(){
            failure:=errors.New("simulated contact read error")
            sql.ExpectQuery(`SELECT .*FROM .*person_contacts.*`).WillReturnError(failure)
            values,err:=contactimpl.NewContactRepository(db).FindByPersonGIDs(context.Background(),[]string{ids.GID})
            Expect(values).To(BeNil())
            Expect(err).To(MatchError(failure))
        })
        It("returns empty IDs without a query for empty contacts",func(){
            gids,err:=contactimpl.NewContactRepository(db).FindPersonGIDsByContacts(context.Background(),nil)
            Expect(err).NotTo(HaveOccurred())
            Expect(gids).To(BeEmpty())
        })
        It("looks up distinct owner IDs for contact values",func(){
            sql.ExpectQuery(`SELECT DISTINCT .*person_gid.*FROM .*person_contacts.*`).
                WillReturnRows(sqlmock.NewRows([]string{"person_gid"}).AddRow(ids.GID))
            gids,err:=contactimpl.NewContactRepository(db).FindPersonGIDsByContacts(context.Background(),[]entity.Contact{fixtures.Contact()})
            Expect(err).NotTo(HaveOccurred())
            Expect(gids).To(Equal([]string{ids.GID}))
        })
        It("propagates an owner lookup error",func(){
            failure:=errors.New("simulated owner lookup error")
            sql.ExpectQuery(`SELECT DISTINCT .*person_gid.*FROM .*person_contacts.*`).WillReturnError(failure)
            _,err:=contactimpl.NewContactRepository(db).FindPersonGIDsByContacts(context.Background(),[]entity.Contact{fixtures.Contact()})
            Expect(err).To(MatchError(failure))
        })
    })
}
