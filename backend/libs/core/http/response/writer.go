package response

import (
	"context"
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	responsedto "grovehub/libs/core/dto/response"
	coreerrors "grovehub/libs/core/errors"
)

func WriteSuccess[T any](writer http.ResponseWriter, statusCode int, data T) error {
	return writeJSON(writer, statusCode, responsedto.NewSuccess(data))
}

func WriteError(writer http.ResponseWriter, statusCode int, definition coreerrors.Definition) error {
	return WriteErrorWithMetadata(writer, statusCode, definition, struct{}{})
}

func WriteErrorContext(
	ctx context.Context,
	writer http.ResponseWriter,
	statusCode int,
	definition coreerrors.Definition,
) error {
	return WriteErrorWithMetadataContext(ctx, writer, statusCode, definition, struct{}{})
}

func WriteErrorWithMetadata[M any](
	writer http.ResponseWriter,
	statusCode int,
	definition coreerrors.Definition,
	metadata M,
) error {
	return writeJSON(
		writer,
		statusCode,
		responsedto.NewError(definition.Code(), definition.Message(), statusCode, metadata),
	)
}

func WriteErrorWithMetadataContext[M any](
	ctx context.Context,
	writer http.ResponseWriter,
	statusCode int,
	definition coreerrors.Definition,
	metadata M,
) error {
	recordErrorOnSpan(ctx, statusCode, definition)
	return WriteErrorWithMetadata(writer, statusCode, definition, metadata)
}

func recordErrorOnSpan(ctx context.Context, statusCode int, definition coreerrors.Definition) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return
	}

	span.SetStatus(codes.Error, definition.Message())
	span.SetAttributes(
		attribute.String("grovehub.error.code", string(definition.Code())),
		attribute.String("grovehub.error.message", definition.Message()),
		attribute.Int("http.response.status_code", statusCode),
	)
	span.AddEvent(
		"grovehub.error.response",
		trace.WithAttributes(
			attribute.String("grovehub.error.code", string(definition.Code())),
			attribute.Int("http.response.status_code", statusCode),
		),
	)
}

func writeJSON[T any](writer http.ResponseWriter, statusCode int, payload T) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(statusCode)
	return json.NewEncoder(writer).Encode(payload)
}
