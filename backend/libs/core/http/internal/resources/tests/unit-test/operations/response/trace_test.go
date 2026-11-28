package response_test

import (
	"context"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/core/http/internal/src/operations/response"
)

var _=Describe("Error response writers",func(){
	It("writes a default error and records a classified error on an active span",func(){
		definition:=coreerrors.NewDefinition("INVALID","Invalid",400)
		w:=httptest.NewRecorder()
		Expect(response.WriteError(w,400,definition)).To(Succeed())
		Expect(w.Body.String()).To(MatchJSON(`{"success":false,"errorCode":"INVALID","errorMessage":"Invalid","httpErrorCode":400,"errorMetaData":{}}`))
		recorder:=tracetest.NewSpanRecorder()
		provider:=sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
		DeferCleanup(func(){Expect(provider.Shutdown(context.Background())).To(Succeed())})
		ctx,span:=provider.Tracer("test").Start(context.Background(),"response")
		Expect(response.WriteErrorContext(ctx,httptest.NewRecorder(),400,definition)).To(Succeed())
		span.End()
		Expect(recorder.Ended()).To(HaveLen(1))
		Expect(recorder.Ended()[0].Status().Code).To(Equal(codes.Error))
		Expect(recorder.Ended()[0].Events()).To(HaveLen(1))
	})
})
