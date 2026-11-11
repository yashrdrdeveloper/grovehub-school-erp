package helpers

import (
	"context"
	"errors"
	"net/http"

	coreerrors "grovehub/libs/core/errors"
	rest "grovehub/libs/core/pipeline/rest"
	httpresponse "grovehub/libs/core/http/response"
)

func WritePersonRequestError(ctx context.Context, writer http.ResponseWriter, err error, invalid coreerrors.Definition) {
	var decodeError *rest.DecodeError
	if errors.As(err, &decodeError) {
		_ = httpresponse.WriteErrorContext(ctx, writer, invalid.HTTPStatus(), invalid)
		return
	}
	WritePersonError(ctx, writer, err, invalid)
}

func WritePersonError(
	ctx context.Context,
	writer http.ResponseWriter,
	err error,
	fallback coreerrors.Definition,
) {
	var applicationError *coreerrors.Error
	if !errors.As(err, &applicationError) {
		_ = httpresponse.WriteErrorContext(ctx, writer, http.StatusInternalServerError, fallback)
		return
	}

	statusCode := applicationError.Definition().HTTPStatus()

	if applicationError.Metadata() == nil {
		_ = httpresponse.WriteErrorContext(ctx, writer, statusCode, applicationError.Definition())
		return
	}

	_ = httpresponse.WriteErrorWithMetadataContext(
		ctx,
		writer,
		statusCode,
		applicationError.Definition(),
		applicationError.Metadata(),
	)
}
