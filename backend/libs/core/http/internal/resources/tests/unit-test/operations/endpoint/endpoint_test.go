package endpoint_test

import (
	"context"
	"net/http"
	"strings"
	corerequest "grovehub/libs/core/request"
	"grovehub/libs/core/http/internal/src/operations/endpoint"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Call", func() {
	It("dispatches through the handler and rejects invalid dependencies and paths", func() {
		type params struct { ID string `path:"id"` }
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Expect(r.URL.Path).To(Equal("/items/a"))
			_, err := w.Write([]byte(`{"success":true,"data":"ok"}`))
			Expect(err).NotTo(HaveOccurred())
		})
		result, status, err := endpoint.Call[string](context.Background(), handler, "GET", "/items/{id}", params{ID:"a"}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal("ok"))
		Expect(status).To(Equal(200))
		_, _, err = endpoint.Call[string](nil, handler, "GET", "/", nil, nil)
		Expect(err).To(HaveOccurred())
		_, _, err = endpoint.Call[string](context.Background(), nil, "GET", "/", nil, nil)
		Expect(err).To(HaveOccurred())
		_, _, err = endpoint.Call[string](context.Background(), handler, "GET", "/items/{id}", params{ID:".."}, nil)
		Expect(err).To(HaveOccurred())
	})

	It("limits the response buffer and rejects malformed envelopes", func() {
		for _, payload := range []string{"not-json", strings.Repeat("x", (8<<20)+1)} {
			handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(payload)) })
			_, _, err := endpoint.Call[string](context.Background(), handler, "GET", "/", nil, nil)
			Expect(err).To(HaveOccurred())
		}
	})
})

var _ = Describe("CallRequest", func() {
	It("encodes declared query values", func() {
		type query struct { Value string `query:"value"` }
		input := endpoint.Request[corerequest.Empty, query, corerequest.Empty]{Query:query{Value:"a b"}}
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Expect(r.URL.Query().Get("value")).To(Equal("a b"))
			w.WriteHeader(204)
		})
		_, status, err := endpoint.CallRequest[string](context.Background(), handler, "GET", "/", input)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(204))
	})
})
