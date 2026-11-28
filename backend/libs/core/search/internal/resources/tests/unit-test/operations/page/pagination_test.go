package page_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math"
	"grovehub/libs/core/search/internal/src/contracts"
	"grovehub/libs/core/search/internal/src/operations/page"
	"grovehub/libs/core/search/internal/src/operations/cursor"
)

var _ = Describe("Pagination",func(){
	It("resolves default, offset and cursor pages",func(){
		for _,item:=range []struct{request contracts.PageRequest;number,rows,offset int}{
			{contracts.PageRequest{},1,25,0},
			{contracts.PageRequest{Rows:50},1,50,0},
			{contracts.PageRequest{PageNumber:3,RowSize:25},3,25,50},
			{contracts.PageRequest{PageNumber:3,RowSize:25,Cursor:"x"},3,25,0},
		}{
			number,rows,offset,err:=page.ResolvePage(item.request,nil)
			Expect(err).NotTo(HaveOccurred())
			Expect([]int{number,rows,offset}).To(Equal([]int{item.number,item.rows,item.offset}))
		}
		for _,request:=range []contracts.PageRequest{{PageNumber:-1},{Rows:10,RowSize:20},{Rows:101},{PageNumber:math.MaxInt},{PageNumber:math.MaxInt/25+2}}{
			_,_,_,err:=page.ResolvePage(request,nil)
			Expect(err).To(MatchError(contracts.ErrInvalidQuery))
		}
	})
	It("binds numbered cursors and calculates result metadata",func(){
		token,err:=cursor.EncodeCursor(contracts.CursorPayload{PageNumber:3,RowSize:25})
		Expect(err).NotTo(HaveOccurred())
		Expect(page.ValidatePageCursor(contracts.PageRequest{},1,25,nil)).To(Succeed())
		Expect(page.ValidatePageCursor(contracts.PageRequest{Cursor:"!"},1,25,nil)).To(MatchError(contracts.ErrInvalidQuery))
		Expect(page.ValidatePageCursor(contracts.PageRequest{PageNumber:3,Cursor:token},3,25,nil)).To(Succeed())
		Expect(page.ValidatePageCursor(contracts.PageRequest{PageNumber:3,Cursor:token},4,25,nil)).To(MatchError(contracts.ErrInvalidQuery))
		Expect(page.NextPageNumber(contracts.PageRequest{PageNumber:3},3)).To(Equal(4))
		Expect(page.NextPageNumber(contracts.PageRequest{},3)).To(BeZero())
		total:=int64(101)
		result:=page.NewPageResponse(contracts.Query[string]{PageNumber:3,Limit:25},contracts.QueryResult[string]{TotalRecords:&total,Items:[]string{"a"},HasNext:true})
		Expect(*result.TotalPages).To(Equal(int64(5)))
		Expect(result.HasPrevious).To(BeTrue())
		Expect(result.HasNext).To(BeTrue())
		Expect(page.NewPageResponse(contracts.Query[string]{},contracts.QueryResult[string]{}).TotalPages).To(BeNil())
	})
})
