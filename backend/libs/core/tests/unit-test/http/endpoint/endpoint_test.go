package endpoint_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreerrors "grovehub/libs/core/errors"
	headercontext "grovehub/libs/core/headercontext"
	endpoint "grovehub/libs/core/http/endpoint"
	httpresponse "grovehub/libs/core/http/response"
)

type personPath struct {
	GID string `path:"gid"`
}

type personRecord struct {
	GID  string `json:"gid"`
	Name string `json:"name"`
}

var _ = Describe("in-process @RequestLine transport", func() {
	It("uses a registered controller route, body, status and typed data", func() {
		mux := http.NewServeMux()
		mux.HandleFunc("POST /foundation/protected/persons", func(writer http.ResponseWriter, request *http.Request) {
			var payload struct { Name string `json:"name"` }
			Expect(request.Method).To(Equal(http.MethodPost))
			Expect(request.Header.Get("Content-Type")).To(Equal("application/json"))
			Expect(request.ContentLength).To(BeNumerically(">", 0))
			decoder := json.NewDecoder(request.Body)
			Expect(decoder.Decode(&payload)).To(Succeed())
			Expect(payload.Name).To(Equal("Ada"))
			Expect(httpresponse.WriteSuccess(writer, http.StatusCreated, personRecord{GID:"gid",Name:payload.Name})).To(Succeed())
		})
		result, status, err := endpoint.Call[personRecord](context.Background(), mux, http.MethodPost, "/foundation/protected/persons", nil, struct{ Name string `json:"name"` }{Name:"Ada"})
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusCreated))
		Expect(result).To(Equal(personRecord{GID:"gid",Name:"Ada"}))
	})

	It("expands only tagged placeholders and preserves caller context and trusted metadata", func() {
		mux := http.NewServeMux()
		typedKey := struct{ Name string }{"request"}
		mux.HandleFunc("GET /foundation/protected/persons/{gid}", func(writer http.ResponseWriter, request *http.Request) {
			Expect(request.PathValue("gid")).To(Equal("valid-gid"))
			Expect(request.Context().Value(typedKey)).To(Equal("value"))
			Expect(headercontext.IdempotencyKey(request.Context())).To(Equal("valid-idempotency-key"))
			Expect(request.Header.Get("Idempotency-Key")).To(Equal("valid-idempotency-key"))
			Expect(httpresponse.WriteSuccess(writer, http.StatusOK, personRecord{GID:request.PathValue("gid")})).To(Succeed())
		})
		ctx := headercontext.WithValues(context.WithValue(context.Background(), typedKey, "value"), headercontext.Values{IdempotencyKey:"valid-idempotency-key"})
		result, status, err := endpoint.Call[personRecord](ctx, mux, http.MethodGet, "/foundation/protected/persons/{gid}", personPath{GID:"valid-gid"}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(status).To(Equal(http.StatusOK))
		Expect(result.GID).To(Equal("valid-gid"))
	})

	It("round-trips semantic error codes, statuses and public metadata", func() {
		mux := http.NewServeMux()
		conflict := coreerrors.NewDefinition("TEST_CONTACT_CONFLICT", "Contact conflicts", http.StatusConflict)
		mux.HandleFunc("PATCH /foundation/protected/persons/{gid}", func(w http.ResponseWriter, r *http.Request) {
			Expect(httpresponse.WriteErrorWithMetadataContext(r.Context(), w, http.StatusConflict, conflict, map[string]string{"gid":r.PathValue("gid")})).To(Succeed())
		})
		_, status, err := endpoint.Call[personRecord](context.Background(), mux, http.MethodPatch, "/foundation/protected/persons/{gid}", personPath{GID:"valid-gid"}, struct{}{})
		Expect(status).To(Equal(http.StatusConflict))
		var applicationError *coreerrors.Error
		Expect(errors.As(err, &applicationError)).To(BeTrue())
		Expect(applicationError.Definition().Code()).To(Equal(conflict.Code()))
		Expect(applicationError.Definition().HTTPStatus()).To(Equal(http.StatusConflict))
		Expect(applicationError.Metadata()).To(Equal(map[string]any{"gid":"valid-gid"}))
	})

	It("preserves nil error metadata from the standard response envelope", func() {
		mux := http.NewServeMux()
		definition := coreerrors.NewDefinition("TEST_NOT_FOUND", "No record", http.StatusNotFound)
		mux.HandleFunc("GET /record", func(writer http.ResponseWriter, request *http.Request) {
			Expect(httpresponse.WriteErrorContext(request.Context(), writer, http.StatusNotFound, definition)).To(Succeed())
		})
		_, status, err := endpoint.Call[personRecord](context.Background(), mux, http.MethodGet, "/record", nil, nil)
		Expect(status).To(Equal(http.StatusNotFound))
		var applicationError *coreerrors.Error
		Expect(errors.As(err, &applicationError)).To(BeTrue())
		Expect(applicationError.Definition().Code()).To(Equal(definition.Code()))
		Expect(applicationError.Metadata()).To(BeNil())
	})

	It("bounds response buffering without accepting truncated JSON", func() {
		handler := http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			_, _ = writer.Write([]byte(strings.Repeat("x", 8<<20+1)))
		})
		_, _, err := endpoint.Call[personRecord](context.Background(), handler, http.MethodGet, "/large", nil, nil)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("exceeds maximum size"))
	})

	It("fails closed for malformed envelopes or unmatched controller routes", func() {
		for _, handler := range []http.Handler{
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "not JSON") }),
			http.NewServeMux(),
		} {
			_, _, err := endpoint.Call[personRecord](context.Background(), handler, http.MethodGet, "/not-registered", nil, nil)
			Expect(err).To(HaveOccurred())
			var appErr *coreerrors.Error
			Expect(errors.As(err, &appErr)).To(BeFalse())
		}
	})

	It("rejects missing, unsafe or unresolved path parameters before dispatch", func() {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /person/{gid}", func(http.ResponseWriter,*http.Request) { Fail("should not dispatch") })
		for _, path := range []personPath{{GID:""},{GID:".."},{GID:"."}} {
			_, _, err := endpoint.Call[personRecord](context.Background(), mux, http.MethodGet, "/person/{gid}", path, nil)
			Expect(err).To(HaveOccurred())
		}
		_, _, err := endpoint.Call[personRecord](context.Background(), mux, http.MethodGet, "/person/{other}", personPath{GID:"ok"}, nil)
		Expect(err).To(HaveOccurred())
	})

	It("handles cancellation before dispatch and limits malformed metadata to safe errors", func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, _, err := endpoint.Call[personRecord](ctx, http.HandlerFunc(func(http.ResponseWriter,*http.Request) { Fail("should not run") }), http.MethodGet, "/anything", nil, nil)
		Expect(errors.Is(err, context.Canceled)).To(BeTrue())

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"success":false,"errorCode":"INVALID","errorMessage":"Invalid","httpErrorCode":404}`)
		})
		_, status, err := endpoint.Call[personRecord](context.Background(), handler, http.MethodGet, "/bad", nil, nil)
		Expect(status).To(Equal(http.StatusBadRequest))
		Expect(err).To(HaveOccurred())
		Expect(strings.Contains(err.Error(), "invalid endpoint error envelope")).To(BeTrue())
	})
})
