package transaction_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
)

var _ = Describe("Public database constructors",func(){
	It("rejects invalid connection configuration",func(){
		db,err:=coredatabase.NewDB("://invalid")
		Expect(err).To(HaveOccurred())
		Expect(db).To(BeNil())
		pool,err:=coredatabase.OpenProjectPool(context.Background(),"postgres://school:secret@localhost/school?user=other")
		Expect(err).To(HaveOccurred())
		Expect(pool).To(BeNil())
		Expect(err.Error()).NotTo(ContainSubstring("secret"))
	})
})
