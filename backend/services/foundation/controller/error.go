package controller

import (
	"context"
	"errors"
	"net/http"

	coreerrors "grovehub/libs/core/errors"
	httpresponse "grovehub/libs/core/http/response"
	personerror "grovehub/services/foundation/constants/error"
)

func writePersonError(
	ctx context.Context,
	writer http.ResponseWriter,
	err error,
	fallback coreerrors.Definition,
) {
	var applicationError *coreerrors.Error
	if !errors.As(err, &applicationError) {
		_ = httpresponse.WriteErrorContext(
			ctx,
			writer,
			http.StatusInternalServerError,
			fallback,
		)
		return
	}

	statusCode := personErrorHTTPStatus(applicationError.Definition())
	metadata := applicationError.Metadata()
	if metadata == nil {
		_ = httpresponse.WriteErrorContext(
			ctx,
			writer,
			statusCode,
			applicationError.Definition(),
		)
		return
	}

	_ = httpresponse.WriteErrorWithMetadataContext(
		ctx,
		writer,
		statusCode,
		applicationError.Definition(),
		metadata,
	)
}

func personErrorHTTPStatus(definition coreerrors.Definition) int {
	switch definition.Code() {
	case personerror.InvalidPersonGID.Code():
		return http.StatusBadRequest
	case personerror.PersonNotFound.Code():
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
