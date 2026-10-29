package helpers

import (
	"context"
	"errors"
	"net/http"

	coreerrors "grovehub/libs/core/errors"
	httpresponse "grovehub/libs/core/http/response"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
)

var errorStatus = map[coreerrors.Code]int{
	mailprovidererror.InvalidDomain.Code():                 http.StatusBadRequest,
	mailprovidererror.DomainAlreadyExists.Code():           http.StatusConflict,
	mailprovidererror.DomainNotFound.Code():                http.StatusNotFound,
	mailprovidererror.DomainNotVerified.Code():             http.StatusConflict,
	mailprovidererror.DomainConcurrentModification.Code():  http.StatusConflict,
	mailprovidererror.DomainCreateFailed.Code():            http.StatusInternalServerError,
	mailprovidererror.DomainReadFailed.Code():              http.StatusInternalServerError,
	mailprovidererror.DomainUpdateFailed.Code():            http.StatusInternalServerError,
	mailprovidererror.DomainVerificationFailed.Code():      http.StatusBadGateway,
	mailprovidererror.InvalidSender.Code():                 http.StatusBadRequest,
	mailprovidererror.SenderAlreadyExists.Code():           http.StatusConflict,
	mailprovidererror.SenderNotFound.Code():                http.StatusNotFound,
	mailprovidererror.SenderConcurrentModification.Code():  http.StatusConflict,
	mailprovidererror.SenderCreateFailed.Code():            http.StatusInternalServerError,
	mailprovidererror.SenderReadFailed.Code():              http.StatusInternalServerError,
	mailprovidererror.SenderUpdateFailed.Code():            http.StatusInternalServerError,
	mailprovidererror.InvalidMessage.Code():                http.StatusBadRequest,
	mailprovidererror.InvalidMessageIdempotencyKey.Code():  http.StatusBadRequest,
	mailprovidererror.MessageIdempotencyConflict.Code():    http.StatusConflict,
	mailprovidererror.MessageSubmissionInProgress.Code():   http.StatusConflict,
	mailprovidererror.MessageRetryNotAllowed.Code():        http.StatusConflict,
	mailprovidererror.MessageResendNotAllowed.Code():       http.StatusConflict,
	mailprovidererror.MessageReplyNotAllowed.Code():        http.StatusConflict,
	mailprovidererror.MessageConcurrentModification.Code(): http.StatusConflict,
	mailprovidererror.MessageNotFound.Code():               http.StatusNotFound,
	mailprovidererror.MessageCreateFailed.Code():           http.StatusInternalServerError,
	mailprovidererror.MessageReadFailed.Code():             http.StatusInternalServerError,
	mailprovidererror.MessageSearchFailed.Code():           http.StatusInternalServerError,
	mailprovidererror.MessageUpdateFailed.Code():           http.StatusInternalServerError,
	mailprovidererror.InvalidMessageSearch.Code():          http.StatusBadRequest,
	mailprovidererror.MessageSubmissionFailed.Code():       http.StatusBadGateway,
}

func WriteError(ctx context.Context, writer http.ResponseWriter, err error) {
	var applicationError *coreerrors.Error
	if !errors.As(err, &applicationError) {
		_ = httpresponse.WriteErrorContext(ctx, writer, http.StatusInternalServerError, mailprovidererror.MessageSubmissionFailed)
		return
	}
	statusCode := errorStatus[applicationError.Definition().Code()]
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	_ = httpresponse.WriteErrorContext(ctx, writer, statusCode, applicationError.Definition())
}
