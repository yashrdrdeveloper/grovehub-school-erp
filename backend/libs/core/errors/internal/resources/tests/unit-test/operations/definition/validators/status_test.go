package validators_test

import (
	"net/http"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors/internal/src/operations/definition/validators"
)

var _ = Describe("HTTPStatus", func() {
	It("defaults to HTTP 500 when no status is supplied", func() {
		Expect(validators.HTTPStatus(nil)).To(Equal(http.StatusInternalServerError))
	})

	DescribeTable("accepts the entire 4xx and 5xx declaration range",
		func(status int) {
			Expect(validators.HTTPStatus([]int{status})).To(Equal(status))
		},
		Entry("lower bound", 400),
		Entry("client error", 499),
		Entry("server error", 500),
		Entry("upper bound", 599),
	)

	DescribeTable("rejects invalid statuses",
		func(status int) {
			Expect(func() { validators.HTTPStatus([]int{status}) }).
				To(PanicWith("error definition requires an HTTP 4xx or 5xx status"))
		},
		Entry("negative", -1),
		Entry("zero", 0),
		Entry("success", 200),
		Entry("below lower bound", 399),
		Entry("above upper bound", 600),
	)

	It("rejects multiple status arguments", func() {
		Expect(func() { validators.HTTPStatus([]int{400, 500}) }).
			To(PanicWith("error definition accepts at most one HTTP status"))
	})
})
