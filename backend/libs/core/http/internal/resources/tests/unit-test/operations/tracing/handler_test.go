package tracing_test

import (
	"net/http"
	"net/http/httptest"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/internal/src/operations/tracing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceHandler", func() {
	It("preserves implicit success and the first explicit status", func() {
		for _, status := range []int{200,400,500} {
			w := httptest.NewRecorder()
			next := http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
				if status != 200 { writer.WriteHeader(status); writer.WriteHeader(201) }
				_, err := writer.Write([]byte("body"))
				Expect(err).NotTo(HaveOccurred())
			})
			tracing.ServiceHandler("identity","test",next).ServeHTTP(w,httptest.NewRequest("GET","/",nil))
			Expect(w.Code).To(Equal(status))
			Expect(w.Body.String()).To(Equal("body"))
		}
	})
})
