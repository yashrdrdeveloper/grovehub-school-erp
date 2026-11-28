package cursor_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"encoding/base64"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/operations/cursor"
)

var _ = Describe("Cursor codec",func(){
	It("round-trips all cursor fields without changing the wire format",func(){
		input:=contracts.CursorPayload{Field:"status",Direction:contracts.DirectionAscending,Value:"ACTIVE",CursorId:"id",SecondarySorts:[]contracts.Sort{{Field:"name",Direction:contracts.DirectionDescending}},SecondaryValues:[]string{"Ada"},PageNumber:2,RowSize:25}
		token,err:=cursor.EncodeCursor(input)
		Expect(err).NotTo(HaveOccurred())
		output,err:=cursor.DecodeCursor(token)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(Equal(input))
	})
	It("rejects invalid base64 and JSON",func(){
		for _,token:=range []string{"!",base64.RawURLEncoding.EncodeToString([]byte("not json"))}{
			_,err:=cursor.DecodeCursor(token)
			Expect(err).To(HaveOccurred())
		}
	})
})
