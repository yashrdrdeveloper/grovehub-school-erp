package rest_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreannotation "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation"
	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
	corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
)

type personPath struct {
	GID string `path:"gid" normalize:"trim,uuidv7" validate:"uuidv7"`
}

type personQuery struct {
	Notify bool `query:"notify"`
	Page *int `query:"page" validate:"required,min=1,max=100"`
	Tags []string `query:"tag"`
}

type personBody struct {
	FirstName string `json:"firstName" normalize:"trim" validate:"notblank,maxrunes=10"`
}

type personResponse struct {
	Name string `json:"name"`
}

func routeRequest(method, path, body, gid string) *http.Request {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	if gid != "" { request.SetPathValue("gid", gid) }
	return request
}

var _ = Describe("HTTP REST executor", func() {
	var invalidPath, invalidQuery, invalidBody, fallback coreerrors.Definition

	BeforeEach(func() {
		invalidPath = coreerrors.NewDefinition("TEST_INVALID_PATH", "Invalid path", http.StatusBadRequest)
		invalidQuery = coreerrors.NewDefinition("TEST_INVALID_QUERY", "Invalid query", http.StatusBadRequest)
		invalidBody = coreerrors.NewDefinition("TEST_INVALID_BODY", "Invalid body", http.StatusBadRequest)
		fallback = coreerrors.NewDefinition("TEST_FAILURE", "Operation failed", http.StatusInternalServerError)
	})

	It("binds path, query and body, validates normalized inputs, and passes typed data to the service", func() {
		gid := "018f5d04-315c-75a9-b890-df0d02e787c1"
		request := routeRequest(http.MethodPatch, "/persons/"+gid+"?notify=true&page=2&tag=a&tag=b", `{"firstName":"  Sima  "}`, gid)
		request = request.WithContext(context.WithValue(request.Context(), struct{ Key string }{"trace"}, "kept"))
		writer := httptest.NewRecorder()
		called := 0
		corepipeline.RunREST(writer, request, corepipeline.Config[personPath, personQuery, personBody, personResponse]{
			InvalidPath: invalidPath,
			InvalidQuery: invalidQuery,
			InvalidBody: invalidBody,
			Fallback: fallback,
			Service: func(ctx context.Context, input corepipeline.Input[personPath, personQuery, personBody]) (corepipeline.Result[personResponse], error) {
				called++
				Expect(ctx.Value(struct{ Key string }{"trace"})).To(Equal("kept"))
				Expect(input.Path.GID).To(Equal(gid))
				Expect(input.Query.Notify).To(BeTrue())
				Expect(*input.Query.Page).To(Equal(2))
				Expect(input.Query.Tags).To(Equal([]string{"a","b"}))
				Expect(input.Body.FirstName).To(Equal("Sima"))
				return corepipeline.Result[personResponse]{Data: personResponse{Name: input.Body.FirstName}, StatusCode: http.StatusCreated, Headers: http.Header{"X-Result": []string{"created"}}}, nil
			},
		})
		Expect(called).To(Equal(1))
		Expect(writer.Code).To(Equal(http.StatusCreated))
		Expect(writer.Header().Get("X-Result")).To(Equal("created"))
		Expect(writer.Body.String()).To(ContainSubstring("Sima"))
	})

	It("rejects invalid path UUIDs before executing the service", func() {
		writer := httptest.NewRecorder()
		called := false
		corepipeline.RunREST(writer, routeRequest(http.MethodGet, "/persons/nope", "", "nope"), corepipeline.Config[personPath, corepipeline.Empty, corepipeline.Empty, personResponse]{
			InvalidPath: invalidPath, Fallback: fallback,
			Service: func(_ context.Context, _ corepipeline.Input[personPath, corepipeline.Empty, corepipeline.Empty]) (corepipeline.Result[personResponse], error) {
				called = true
				return corepipeline.OK(personResponse{}), nil
			},
		})
		Expect(called).To(BeFalse())
		Expect(writer.Code).To(Equal(http.StatusBadRequest))
		Expect(writer.Body.String()).To(ContainSubstring("TEST_INVALID_PATH"))
		Expect(writer.Body.String()).NotTo(ContainSubstring("nope"))
	})

	It("rejects invalid or repeated scalar query values and missing required query parameters", func() {
		for _, query := range []string{"?page=0", "?page=101", "?page=NaN", "?page=1&page=2", "?notify=wrong&page=1", ""} {
			writer := httptest.NewRecorder()
			corepipeline.RunREST(writer, routeRequest(http.MethodGet, "/people"+query, "", ""), corepipeline.Config[corepipeline.Empty, personQuery, corepipeline.Empty, personResponse]{
				InvalidQuery: invalidQuery, Fallback: fallback,
				Service: func(_ context.Context, _ corepipeline.Input[corepipeline.Empty, personQuery, corepipeline.Empty]) (corepipeline.Result[personResponse], error) {
					Fail("service must not be called for invalid query")
					return corepipeline.OK(personResponse{}), nil
				},
			})
			Expect(writer.Code).To(Equal(http.StatusBadRequest), query)
			Expect(writer.Body.String()).To(ContainSubstring("TEST_INVALID_QUERY"), query)
		}
	})

	It("validates trimmed body fields before custom validation", func() {
		writer := httptest.NewRecorder()
		customCalled := false
		corepipeline.RunREST(writer, routeRequest(http.MethodPost, "/people", `{"firstName":"  "}`, ""), corepipeline.Config[corepipeline.Empty, corepipeline.Empty, personBody, personResponse]{
			Body: coreannotation.ValidationRules[personBody]{
				Validate: func(_ personBody) error { customCalled = true; return nil },
			},
			InvalidBody: invalidBody, Fallback: fallback,
			Service: func(_ context.Context, _ corepipeline.Input[corepipeline.Empty, corepipeline.Empty, personBody]) (corepipeline.Result[personResponse], error) {
				Fail("service must not be called")
				return corepipeline.OK(personResponse{}), nil
			},
		})
		Expect(customCalled).To(BeFalse())
		Expect(writer.Code).To(Equal(http.StatusBadRequest))
		Expect(writer.Body.String()).To(ContainSubstring("TEST_INVALID_BODY"))
	})

	It("reports strict JSON decode failures with the configured safe error", func() {
		writer := httptest.NewRecorder()
		corepipeline.RunREST(writer, routeRequest(http.MethodPost, "/people", `{"firstName":"Ana","unknown":"secret"}`, ""), corepipeline.Config[corepipeline.Empty, corepipeline.Empty, personBody, personResponse]{
			Decode: corehttp.DecodeStrictJSON,
			InvalidBody: invalidBody,
			Fallback: fallback,
			Service: func(_ context.Context, _ corepipeline.Input[corepipeline.Empty, corepipeline.Empty, personBody]) (corepipeline.Result[personResponse], error) {
				Fail("service must not be called")
				return corepipeline.OK(personResponse{}), nil
			},
		})
		Expect(writer.Code).To(Equal(http.StatusBadRequest))
		Expect(writer.Body.String()).NotTo(ContainSubstring("secret"))
	})

	It("writes the definition status and metadata from service application errors", func() {
		writer := httptest.NewRecorder()
		conflict := coreerrors.NewDefinition("TEST_CONFLICT", "Conflicting data", http.StatusConflict)
		corepipeline.RunREST(writer, routeRequest(http.MethodGet, "/people", "", ""), corepipeline.Config[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty, personResponse]{
			Fallback: fallback,
			Service: func(_ context.Context, _ corepipeline.Input[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty]) (corepipeline.Result[personResponse], error) {
				return corepipeline.Result[personResponse]{}, coreerrors.NewErrorWithMetadata(conflict, map[string]string{"field":"email"})
			},
		})
		Expect(writer.Code).To(Equal(http.StatusConflict))
		Expect(writer.Body.String()).To(ContainSubstring("TEST_CONFLICT"))
		Expect(writer.Body.String()).To(ContainSubstring("email"))
	})

	It("uses a safe HTTP 500 fallback for unexpected service errors", func() {
		writer := httptest.NewRecorder()
		corepipeline.RunREST(writer, routeRequest(http.MethodGet, "/people", "", ""), corepipeline.Config[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty, personResponse]{
			InvalidBody: invalidBody,
			Fallback: fallback,
			Service: func(_ context.Context, _ corepipeline.Input[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty]) (corepipeline.Result[personResponse], error) {
				return corepipeline.Result[personResponse]{}, errors.New("secret database failure")
			},
		})
		Expect(writer.Code).To(Equal(http.StatusInternalServerError))
		Expect(writer.Body.String()).To(ContainSubstring("TEST_FAILURE"))
		Expect(writer.Body.String()).NotTo(ContainSubstring("secret database failure"))
	})

	It("supports HTTP 204 without writing an envelope", func() {
		writer := httptest.NewRecorder()
		corepipeline.RunREST(writer, routeRequest(http.MethodDelete, "/people", "", ""), corepipeline.Config[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty, personResponse]{
			Fallback: fallback,
			Service: func(_ context.Context, _ corepipeline.Input[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty]) (corepipeline.Result[personResponse], error) {
				return corepipeline.Result[personResponse]{StatusCode: http.StatusNoContent}, nil
			},
		})
		Expect(writer.Code).To(Equal(http.StatusNoContent))
		Expect(writer.Body.String()).To(BeEmpty())
	})

	It("supports custom success and error writers without invoking the default writer", func() {
		writer := httptest.NewRecorder()
		called := 0
		corepipeline.RunREST(writer, routeRequest(http.MethodGet, "/people", "", ""), corepipeline.Config[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty, personResponse]{
			Fallback: fallback,
			Writer: corepipeline.Writer[personResponse]{
				Success: func(w http.ResponseWriter, _ *http.Request, result corepipeline.Result[personResponse]) error {
					called++
					w.WriteHeader(http.StatusAccepted)
					return nil
				},
			},
			Service: func(_ context.Context, _ corepipeline.Input[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty]) (corepipeline.Result[personResponse], error) {
				return corepipeline.OK(personResponse{}), nil
			},
		})
		Expect(called).To(Equal(1))
		Expect(writer.Code).To(Equal(http.StatusAccepted))

		writer = httptest.NewRecorder()
		corepipeline.RunREST(writer, routeRequest(http.MethodGet, "/people", "", ""), corepipeline.Config[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty, personResponse]{
			Fallback: fallback,
			Writer: corepipeline.Writer[personResponse]{
				Error: func(w http.ResponseWriter, _ *http.Request, failure corepipeline.Failure) error {
					called++
					Expect(failure.StatusCode).To(Equal(http.StatusInternalServerError))
					w.WriteHeader(http.StatusTeapot)
					return nil
				},
			},
			Service: func(_ context.Context, _ corepipeline.Input[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty]) (corepipeline.Result[personResponse], error) {
				return corepipeline.Result[personResponse]{}, errors.New("unexpected")
			},
		})
		Expect(called).To(Equal(2))
		Expect(writer.Code).To(Equal(http.StatusTeapot))
	})

	It("does not invoke service when a final cross-source validator rejects input", func() {
		writer := httptest.NewRecorder()
		corepipeline.RunREST(writer, routeRequest(http.MethodPost, "/people", `{"firstName":"Alex"}`, ""), corepipeline.Config[corepipeline.Empty, corepipeline.Empty, personBody, personResponse]{
			InvalidBody: invalidBody,
			Fallback: fallback,
			Validate: func(input corepipeline.Input[corepipeline.Empty, corepipeline.Empty, personBody]) error {
				Expect(input.Body.FirstName).To(Equal("Alex"))
				return coreerrors.NewError(invalidBody)
			},
			Service: func(_ context.Context, _ corepipeline.Input[corepipeline.Empty, corepipeline.Empty, personBody]) (corepipeline.Result[personResponse], error) {
				Fail("service should not run")
				return corepipeline.OK(personResponse{}), nil
			},
		})
		Expect(writer.Code).To(Equal(http.StatusBadRequest))
	})

	It("rejects malformed query escaping and unknown keys in strict mode", func() {
		for _, options := range []struct { raw string; strict bool }{
			{"/people?page=1&unknown=%XX", false},
			{"/people?page=1&unexpected=1", true},
		} {
			writer := httptest.NewRecorder()
			corepipeline.RunREST(writer, routeRequest(http.MethodGet, options.raw, "", ""), corepipeline.Config[corepipeline.Empty, personQuery, corepipeline.Empty, personResponse]{
				StrictQuery: options.strict,
				InvalidQuery: invalidQuery,
				Fallback: fallback,
				Service: func(_ context.Context, _ corepipeline.Input[corepipeline.Empty, personQuery, corepipeline.Empty]) (corepipeline.Result[personResponse], error) {
					Fail("service must not run with invalid query")
					return corepipeline.OK(personResponse{}), nil
				},
			})
			Expect(writer.Code).To(Equal(http.StatusBadRequest), options.raw)
		}
	})

	It("writes a safe error envelope for an invalid service result status", func() {
		writer := httptest.NewRecorder()
		corepipeline.RunREST(writer, routeRequest(http.MethodGet, "/people", "", ""), corepipeline.Config[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty, personResponse]{
			Fallback: fallback,
			Service: func(_ context.Context, _ corepipeline.Input[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty]) (corepipeline.Result[personResponse], error) {
				return corepipeline.Result[personResponse]{StatusCode: http.StatusConflict}, nil
			},
		})
		Expect(writer.Code).To(Equal(http.StatusInternalServerError))
		Expect(writer.Body.String()).To(ContainSubstring("TEST_FAILURE"))
	})

	It("defaults an omitted successful Result status to 200", func() {
		writer := httptest.NewRecorder()
		corepipeline.RunREST(writer, routeRequest(http.MethodGet, "/people", "", ""), corepipeline.Config[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty, personResponse]{
			Fallback: fallback,
			Service: func(_ context.Context, _ corepipeline.Input[corepipeline.Empty, corepipeline.Empty, corepipeline.Empty]) (corepipeline.Result[personResponse], error) {
				return corepipeline.Result[personResponse]{Data: personResponse{Name:"OK"}}, nil
			},
		})
		Expect(writer.Code).To(Equal(http.StatusOK))
	})
})
