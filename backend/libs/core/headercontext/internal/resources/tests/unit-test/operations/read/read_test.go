package read_test

import (
	"grovehub/libs/core/headercontext/internal/src/contracts"
	"grovehub/libs/core/headercontext/internal/src/operations/read"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTPHeaders", func() {
	It("returns only allowlisted metadata and accepts absent headers", func() {
		values, err := read.HTTPHeaders(nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(values).To(Equal(contracts.Values{}))
		header := make(http.Header)
		header.Set("Idempotency-Key", " create:123 ")
		header.Set("X-Correlation-ID", " corr-1 ")
		header.Set("X-Request-ID", " req-1 ")
		header.Set("X-Trace-ID", "forged")
		header.Set("Authorization", "untrusted")
		values, err = read.HTTPHeaders(header)
		Expect(err).NotTo(HaveOccurred())
		Expect(values).To(Equal(contracts.Values{IdempotencyKey: "create:123", CorrelationID: "corr-1", RequestID: "req-1"}))
	})

	It("rejects invalid or duplicate metadata without returning partial values", func() {
		for _, name := range []string{"Idempotency-Key", "X-Correlation-ID", "X-Request-ID"} {
			for _, inputs := range [][]string{{" "}, {"bad value"}, {"first", "second"}} {
				header := make(http.Header)
				header.Set("Idempotency-Key", "valid")
				header.Set("X-Correlation-ID", "valid")
				header.Set("X-Request-ID", "valid")
				header.Del(name)
				for _, input := range inputs {
					header.Add(name, input)
				}
				values, err := read.HTTPHeaders(header)
				Expect(err).To(BeIdenticalTo(contracts.ErrInvalidHeader))
				Expect(values).To(Equal(contracts.Values{}))
			}
		}
		for _, name := range []string{"X-Correlation-ID", "X-Request-ID"} {
			header := make(http.Header)
			header.Set(name, strings.Repeat("x", 129))
			_, err := read.HTTPHeaders(header)
			Expect(err).To(BeIdenticalTo(contracts.ErrInvalidHeader))
		}
	})
})
