package impl

import (
	coregateway "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	headercontext "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext"
	coretelemetry "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	observability "grovehub/launcher/internal/src/middleware/observability/impl"
	"log/slog"
	"net/http"
	"strings"
)

func newEdgeHandler(logger *slog.Logger, mux *http.ServeMux, external http.Handler) http.Handler {
	restrictedHandler := coregateway.RestrictedOnly(mux, external)
	edgeService := coretelemetry.ServiceLauncherEdge
	edgeHandler := otelhttp.NewHandler(
		observability.New(logger).Wrap(
			headercontext.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, pattern := mux.Handler(r)
				if pattern == "GET /health" || pattern == "GET /ready" {
					external.ServeHTTP(w, r)
					return
				}

				restrictedHandler.ServeHTTP(w, r)
			})),
		),
		"grovehub.edge",
		otelhttp.WithTracerProvider(coretelemetry.TracerProvider(edgeService.Key)),
		otelhttp.WithMeterProvider(coretelemetry.MeterProvider(edgeService.Key)),
		otelhttp.WithSpanNameFormatter(func(_ string, request *http.Request) string {
			route := strings.TrimPrefix(request.Pattern, request.Method+" ")
			if route == "" {
				route = "unmatched"
			}

			return edgeService.Name + " " + request.Method + " " + route
		}),
		otelhttp.WithSpanOptions(
			trace.WithAttributes(
				attribute.String("grovehub.service.key", edgeService.Key),
				attribute.String("grovehub.service.display_name", edgeService.Name),
			),
		),
	)

	return edgeHandler
}
