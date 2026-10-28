package helpers

import (
	"context"
	"errors"
	"net/http"

	coreerrors "grovehub/libs/core/errors"
	httpresponse "grovehub/libs/core/http/response"
	identityerror "grovehub/services/identity/constants/error"
)

var registerErrorStatus = map[coreerrors.Code]int{
	identityerror.InvalidRegistration.Code():               http.StatusBadRequest,
	identityerror.InvalidRegistrationIdempotencyKey.Code(): http.StatusBadRequest,
	identityerror.RegistrationAlreadyExists.Code():         http.StatusConflict,
	identityerror.RegistrationIdempotencyConflict.Code():   http.StatusConflict,
	identityerror.RegistrationFailed.Code():                http.StatusInternalServerError,
}

func WriteRegisterError(ctx context.Context, writer http.ResponseWriter, err error) {
	var applicationError *coreerrors.Error
	if !errors.As(err, &applicationError) {
		_ = httpresponse.WriteErrorContext(
			ctx,
			writer,
			http.StatusInternalServerError,
			identityerror.RegistrationFailed,
		)
		return
	}

	statusCode := registerErrorStatus[applicationError.Definition().Code()]
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	_ = httpresponse.WriteErrorContext(ctx, writer, statusCode, applicationError.Definition())
}
