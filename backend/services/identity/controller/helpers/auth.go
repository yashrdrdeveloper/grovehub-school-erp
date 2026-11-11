package helpers

import (
	"context"
	"errors"
	"net/http"

	coreerrors "grovehub/libs/core/errors"
	httpresponse "grovehub/libs/core/http/response"
	identityerror "grovehub/services/identity/constants/error"
)

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

	statusCode := applicationError.Definition().HTTPStatus()
	_ = httpresponse.WriteErrorContext(ctx, writer, statusCode, applicationError.Definition())
}
