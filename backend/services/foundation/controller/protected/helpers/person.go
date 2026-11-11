package helpers

import (
	"context"
	"errors"
	"net/http"

	coreerrors "grovehub/libs/core/errors"
	rest "grovehub/libs/core/pipeline/rest"
	httpresponse "grovehub/libs/core/http/response"
	personerror "grovehub/services/foundation/constants/error"
)

var personErrorStatus = map[coreerrors.Code]int{
	personerror.InvalidPersonGID.Code():             http.StatusBadRequest,
	personerror.InvalidPersonRequest.Code():         http.StatusBadRequest,
	personerror.InvalidPersonContact.Code():         http.StatusBadRequest,
	personerror.InvalidPersonIdempotencyKey.Code():  http.StatusBadRequest,
	personerror.InvalidPersonSearch.Code():          http.StatusBadRequest,
	personerror.InvalidPersonSearchCursor.Code():    http.StatusBadRequest,
	personerror.PersonNotFound.Code():               http.StatusNotFound,
	personerror.PersonContactConflict.Code():        http.StatusConflict,
	personerror.PersonIdempotencyConflict.Code():    http.StatusConflict,
	personerror.PersonConcurrentModification.Code(): http.StatusConflict,
}

func WritePersonRequestError(ctx context.Context, writer http.ResponseWriter, err error, invalid coreerrors.Definition) {
	var decodeError *rest.DecodeError
	if errors.As(err, &decodeError) {
		_ = httpresponse.WriteErrorContext(ctx, writer, http.StatusBadRequest, invalid)
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

	statusCode := personErrorStatus[applicationError.Definition().Code()]
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

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
