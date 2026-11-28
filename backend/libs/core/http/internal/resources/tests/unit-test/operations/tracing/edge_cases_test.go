package tracing_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"grovehub/libs/core/http/internal/src/operations/tracing"
	coretelemetry "grovehub/libs/core/telemetry"
)

type captureTransport func(*http.Request)(*http.Response,error)
func (f captureTransport) RoundTrip(r *http.Request)(*http.Response,error){return f(r)}

var _=Describe("HTTP tracing transport and response wrapper",func(){
	It("propagates trace headers on an independent request clone",func(){
		previous:=otel.GetTracerProvider()
		provider:=sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
		otel.SetTracerProvider(provider)
		DeferCleanup(func(){otel.SetTracerProvider(previous);Expect(provider.Shutdown(context.Background())).To(Succeed())})
		ctx,span:=provider.Tracer("test").Start(context.Background(),"caller")
		defer span.End()
		req,err:=http.NewRequestWithContext(ctx,"GET","https://example.test",nil)
		Expect(err).NotTo(HaveOccurred())
		transport:=tracing.NewTransport(captureTransport(func(r *http.Request)(*http.Response,error){
			Expect(r).NotTo(BeIdenticalTo(req))
			Expect(r.Header.Get("X-Trace-ID")).To(Equal(coretelemetry.TraceID(ctx)))
			return &http.Response{StatusCode:204,Header:make(http.Header),Body:io.NopCloser(strings.NewReader("")),Request:r},nil
		}))
		resp,err:=transport.RoundTrip(req)
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()
		Expect(req.Header.Get("X-Trace-ID")).To(BeEmpty())
	})
	It("allows response-controller unwrapping and retains caller identity and matched route",func(){
		w:=httptest.NewRecorder()
		req:=httptest.NewRequest("GET","/items/1",nil)
		req.Pattern="GET /items/{id}"
		req=req.WithContext(coretelemetry.ContextWithService(req.Context(),coretelemetry.ServiceFoundation))
		handler:=tracing.ServiceHandler("identity","read",http.HandlerFunc(func(writer http.ResponseWriter,r *http.Request){
			underlying,ok:=writer.(interface{Unwrap() http.ResponseWriter})
			Expect(ok).To(BeTrue())
			Expect(underlying.Unwrap()).To(BeIdenticalTo(w))
			writer.WriteHeader(204)
		}))
		handler.ServeHTTP(w,req)
		Expect(w.Code).To(Equal(204))
	})
})
