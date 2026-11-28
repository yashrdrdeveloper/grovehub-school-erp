package contracts_test

import (
	"net/http"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/module/internal/src/contracts"
)

var _ = Describe("Module", func() {
	It("defines only HTTP registration and shutdown", func() {
		contract := reflect.TypeFor[contracts.Module]()
		Expect(contract.Kind()).To(Equal(reflect.Interface))
		Expect(contract.NumMethod()).To(Equal(2))

		register, ok := contract.MethodByName("RegisterHTTP")
		Expect(ok).To(BeTrue())
		Expect(register.Type).To(Equal(reflect.TypeFor[func(*http.ServeMux)]()))

		closeMethod, ok := contract.MethodByName("Close")
		Expect(ok).To(BeTrue())
		Expect(closeMethod.Type).To(Equal(reflect.TypeFor[func()]()))
	})
})
