package tracing_test

import (
	"grovehub/libs/core/http/internal/src/operations/tracing"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewClient", func() {
	It("preserves timeout and wraps the default transport", func() {
		client := tracing.NewClient(time.Second)
		Expect(client.Timeout).To(Equal(time.Second))
		Expect(client.Transport).NotTo(BeNil())
		Expect(tracing.NewTransport(nil)).NotTo(BeNil())
		Expect(tracing.NewServiceTransport("identity", http.DefaultTransport)).NotTo(BeNil())
		Expect(tracing.NewServiceClient("identity", 2*time.Second).Timeout).To(Equal(2 * time.Second))
	})
})
