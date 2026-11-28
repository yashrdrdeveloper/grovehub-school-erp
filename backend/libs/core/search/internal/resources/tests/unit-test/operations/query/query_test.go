package query_test

import (
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/operations/query"
)

type record struct {
	ID string
	Name string
	Status string
}

func newDB(dryRun bool) (*gorm.DB, sqlmock.Sqlmock) {
	connection, mock, err := sqlmock.New()
	Expect(err).NotTo(HaveOccurred())
	DeferCleanup(func(){ _ = connection.Close() })
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: connection}), &gorm.Config{DryRun:dryRun, DisableAutomaticPing:true})
	Expect(err).NotTo(HaveOccurred())
	return db,mock
}

func queryConfig() contracts.Config {
	return contracts.Config{
		SearchExpression:"name", CursorExpression:"id",
		ScalarFields:map[string]contracts.ScalarField{"status":{Expression:"status"}},
		SortFields:map[string]contracts.SortField{"id":{Expression:"id"},"name":{Expression:"name"},"status":{Expression:"status"}},
	}
}

func baseQuery() contracts.Query[string] {
	return contracts.Query[string]{Sort:contracts.QuerySort{Field:"name",Direction:contracts.DirectionAscending},Limit:2}
}

var _ = Describe("Apply",func(){
	It("creates parameterized search, filter and pagination clauses",func(){
		db,_:=newDB(true)
		input:=baseQuery()
		input.SearchTerms=[]string{"a%b"}
		input.Filters=[]contracts.QueryFilter{{Field:"status",Operator:contracts.OperatorEquals,Values:[]string{"ACTIVE"}}}
		input.Offset=4
		result,err:=query.Apply(db.Model(&record{}),input,queryConfig())
		Expect(err).NotTo(HaveOccurred())
		statement:=result.Find(&[]record{}).Statement
		Expect(statement.SQL.String()).To(ContainSubstring("name ILIKE"))
		Expect(statement.SQL.String()).To(ContainSubstring("ORDER BY name ASC,id ASC"))
		Expect(statement.SQL.String()).To(ContainSubstring("LIMIT"))
		Expect(statement.SQL.String()).To(ContainSubstring("OFFSET"))
		Expect(statement.Vars).To(ContainElement(`%a\%b%`))
		Expect(statement.Vars).To(ContainElement("ACTIVE"))
	})
	It("rejects nil DB, invalid query and unknown filters",func(){
		db,_:=newDB(true)
		_,err:=query.Apply(nil,baseQuery(),queryConfig())
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
		input:=baseQuery();input.Limit=0
		_,err=query.Apply(db,input,queryConfig())
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
		input=baseQuery();input.Filters=[]contracts.QueryFilter{{Field:"unknown"}}
		_,err=query.Apply(db,input,queryConfig())
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
	})
	It("uses keyset conditions and the cursor expression as a stable tie-breaker",func(){
		db,_:=newDB(true)
		input:=baseQuery()
		input.Cursor=&contracts.Cursor[string]{Value:"Ada",CursorId:"id-1"}
		result,err:=query.Apply(db.Model(&record{}),input,queryConfig())
		Expect(err).NotTo(HaveOccurred())
		statement:=result.Find(&[]record{}).Statement
		Expect(statement.SQL.String()).To(ContainSubstring("name >"))
		Expect(statement.SQL.String()).To(ContainSubstring("id >"))
		Expect(statement.Vars).To(ContainElement("id-1"))
	})
})

var _ = Describe("Execute",func(){
	It("fetches one extra row and returns only the requested page",func(){
		db,mock:=newDB(false)
		mock.ExpectQuery("(?i)SELECT .*FROM").WillReturnRows(sqlmock.NewRows([]string{"id","name"}).AddRow("1","A").AddRow("2","B").AddRow("3","C"))
		result,err:=query.Execute[string,record](db,baseQuery(),queryConfig())
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Items).To(HaveLen(2))
		Expect(result.HasNext).To(BeTrue())
		Expect(result.TotalRecords).To(BeNil())
		Expect(mock.ExpectationsWereMet()).To(Succeed())
	})
	It("counts numbered pages separately and preserves total records",func(){
		db,mock:=newDB(false)
		input:=baseQuery();input.PageNumber=2;input.Offset=2
		mock.ExpectQuery("(?i)SELECT count").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
		mock.ExpectQuery("(?i)SELECT .*FROM").WillReturnRows(sqlmock.NewRows([]string{"id","name"}).AddRow("3","C"))
		result,err:=query.Execute[string,record](db,input,queryConfig())
		Expect(err).NotTo(HaveOccurred())
		Expect(*result.TotalRecords).To(Equal(int64(3)))
		Expect(result.Items).To(HaveLen(1))
		Expect(result.HasNext).To(BeFalse())
		Expect(mock.ExpectationsWereMet()).To(Succeed())
	})
	It("propagates database read and count failures",func(){
		for _,number:=range []int{0,1}{
			db,mock:=newDB(false)
			input:=baseQuery();input.PageNumber=number
			sentinel:=errors.New("database unavailable")
			mock.ExpectQuery("(?i)SELECT").WillReturnError(sentinel)
			_,err:=query.Execute[string,record](db,input,queryConfig())
			Expect(err).To(MatchError(sentinel))
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		}
	})
})
