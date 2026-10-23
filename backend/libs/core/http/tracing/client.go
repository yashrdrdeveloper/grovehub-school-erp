package tracing

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	httpheaders "grovehub/libs/core/http/headers"
	coretelemetry "grovehub/libs/core/telemetry"
)

type traceHeaderTransport struct {
	next http.RoundTripper
}

func NewTransport(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}

	return otelhttp.NewTransport(&traceHeaderTransport{
		next: base,
	})
}

func NewClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: NewTransport(http.DefaultTransport),
		Timeout:   timeout,
	}
}

func (transport *traceHeaderTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	clonedRequest := request.Clone(request.Context())
	if traceID := coretelemetry.TraceID(request.Context()); traceID != "" {
		clonedRequest.Header.Set(httpheaders.TraceID, traceID)
	}

	return transport.next.RoundTrip(clonedRequest)
}
