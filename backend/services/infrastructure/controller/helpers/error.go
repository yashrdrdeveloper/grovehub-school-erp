package helpers

import (
	"context"
	"errors"
	"net/http"

	coreerrors "grovehub/libs/core/errors"
	httpresponse "grovehub/libs/core/http/response"
	infrastructureerror "grovehub/services/infrastructure/constants/error"
)

var errorStatus = map[coreerrors.Code]int{
	infrastructureerror.InvalidDKIMIdentity.Code():  http.StatusBadRequest,
	infrastructureerror.DKIMSelectorConflict.Code(): http.StatusConflict,
	infrastructureerror.DKIMPrepareFailed.Code():     http.StatusInternalServerError,
}

func WriteError(ctx context.Context, writer http.ResponseWriter, err error) {
	var applicationError *coreerrors.Error
	if !errors.As(err, &applicationError) {
		_ = httpresponse.WriteErrorContext(
			ctx,
			writer,
			http.StatusInternalServerError,
			infrastructureerror.DKIMPrepareFailed,
		)
		return
	}
	statusCode := errorStatus[applicationError.Definition().Code()]
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	_ = httpresponse.WriteErrorContext(ctx, writer, statusCode, applicationError.Definition())
}
