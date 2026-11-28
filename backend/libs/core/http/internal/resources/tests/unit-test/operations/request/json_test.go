package request_test

import (
	"net/http/httptest"
	"strings"
	"grovehub/libs/core/http/internal/src/operations/request"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DecodeJSON", func() {
	It("accepts a single value and rejects malformed or trailing JSON", func() {
		for _, raw := range []string{`{"name":"a"}`,`{"name":"a","extra":true}`} {
			var value struct { Name string `json:"name"` }
			Expect(request.DecodeJSON(httptest.NewRequest("POST","/",strings.NewReader(raw)), &value)).To(Succeed())
			Expect(value.Name).To(Equal("a"))
		}
		for _, raw := range []string{"", "{", "{} {}", "{} invalid"} {
			var value map[string]any
			Expect(request.DecodeJSON(httptest.NewRequest("POST","/",strings.NewReader(raw)), &value)).NotTo(Succeed())
		}
	})
})

var _ = Describe("DecodeStrictJSON", func() {
	It("also rejects unknown fields", func() {
		for _, raw := range []string{`{"extra":true}`, "{} {}", "{"} {
			var value struct { Name string `json:"name"` }
			Expect(request.DecodeStrictJSON(httptest.NewRequest("POST","/",strings.NewReader(raw)), &value)).NotTo(Succeed())
		}
		var value struct { Name string `json:"name"` }
		Expect(request.DecodeStrictJSON(httptest.NewRequest("POST","/",strings.NewReader(`{"name":"a"}`)), &value)).To(Succeed())
		Expect(value.Name).To(Equal("a"))
	})
})
