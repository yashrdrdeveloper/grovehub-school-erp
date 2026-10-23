package response_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	httpresponse "grovehub/libs/core/http/response"
)

var _ = Describe("HTTP response writer", func() {
	It("writes a semantic error definition with the HTTP status", func() {
		recorder := httptest.NewRecorder()
		definition := serviceDownDefinition()

		Expect(httpresponse.WriteError(recorder, http.StatusServiceUnavailable, definition)).To(Succeed())
		payload := decodeDefaultError(recorder)

		Expect(recorder.Code).To(Equal(http.StatusServiceUnavailable))
		Expect(payload.Success).To(BeFalse())
		Expect(payload.ErrorCode).To(Equal(definition.Code()))
		Expect(payload.ErrorMessage).To(Equal(definition.Message()))
		Expect(payload.HTTPErrorCode).To(Equal(http.StatusServiceUnavailable))
		Expect(payload.ErrorMetaData).To(BeEmpty())
	})

	It("writes typed error metadata", func() {
		recorder := httptest.NewRecorder()

		Expect(httpresponse.WriteErrorWithMetadata(
			recorder,
			http.StatusBadRequest,
			invalidRequestDefinition(),
			metadata{Field: metadataField},
		)).To(Succeed())

		Expect(decodeMetadataError(recorder).ErrorMetaData.Field).To(Equal(metadataField))
	})

	It("annotates the active span with the semantic error definition", func() {
		ctx, span := newErrorTraceContext()
		recorder := httptest.NewRecorder()
		definition := serviceDownDefinition()

		Expect(httpresponse.WriteErrorContext(
			ctx,
			recorder,
			http.StatusServiceUnavailable,
			definition,
		)).To(Succeed())
		span.End()

		endedSpans := responseSpanRecorder.Ended()
		Expect(endedSpans).To(HaveLen(1))
		Expect(endedSpans[0].Status().Code).To(Equal(codes.Error))
		Expect(endedSpans[0].Attributes()).To(ContainElement(
			attribute.String("grovehub.error.code", serviceDownCode),
		))
		Expect(endedSpans[0].Events()).NotTo(BeEmpty())
	})
})
