package validation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"errors"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/operations/cursor"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/operations/validation"
)

var _ = Describe("Cursor validation",func(){
	It("validates UUID, time and non-empty cursor values",func(){
		sentinel:=errors.New("invalid")
		Expect(validation.ValidateUUIDv7Cursor("019ab812-3456-7890-8123-123456789abc",nil)).To(Succeed())
		Expect(validation.ValidateUUIDv7Cursor("bad",sentinel)).To(MatchError(sentinel))
		Expect(validation.ValidateTimeCursor("2026-11-18T00:00:00Z",nil)).To(Succeed())
		Expect(validation.ValidateTimeCursor("bad",nil)).To(MatchError(contracts.ErrInvalidQuery))
		Expect(validation.ValidateNonEmptyCursor("value",nil)).To(Succeed())
		Expect(validation.ValidateNonEmptyCursor(" ",nil)).To(MatchError(contracts.ErrInvalidQuery))
	})
	It("binds a decoded cursor to the requested sort and validates values",func(){
		sort:=contracts.QuerySort{Field:"name",Direction:contracts.DirectionAscending}
		token,err:=cursor.EncodeCursor(contracts.CursorPayload{Field:sort.Field,Direction:sort.Direction,Value:"Ada",CursorId:"id"})
		Expect(err).NotTo(HaveOccurred())
		result,_,err:=validation.NormalizeCursor("",sort,nil,nil,nil)
		Expect(err).NotTo(HaveOccurred());Expect(result).To(BeNil())
		result,_,err=validation.NormalizeCursor(token,sort,validation.ValidateNonEmptyCursor,validation.ValidateNonEmptyCursor,nil)
		Expect(err).NotTo(HaveOccurred());Expect(result.CursorId).To(Equal("id"))
		_,_,err=validation.NormalizeCursor(token,sort,nil,nil,nil);Expect(err).To(HaveOccurred())
		sort.Direction=contracts.DirectionDescending
		_,_,err=validation.NormalizeCursor(token,sort,validation.ValidateNonEmptyCursor,validation.ValidateNonEmptyCursor,nil);Expect(err).To(HaveOccurred())
	})
})
