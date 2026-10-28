package impl

import (
	"strings"

	coreerrors "grovehub/libs/core/errors"
	coreidempotency "grovehub/libs/core/idempotency"
	coresearch "grovehub/libs/core/search"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
	"grovehub/services/mailProvider/entity"
	mailproviderservices "grovehub/services/mailProvider/services"
	mailprovidervalidators "grovehub/services/mailProvider/validators"
	validatorhelpers "grovehub/services/mailProvider/validators/helpers"
)

type MessageValidator struct{}

func NewMessageValidator() mailprovidervalidators.MessageValidator { return &MessageValidator{} }

func (validator *MessageValidator) ValidateSend(
	input mailproviderservices.SendMessageInput,
) (mailproviderservices.SendMessageInput, error) {
	idempotencyKey, err := coreidempotency.ValidateKey(input.IdempotencyKey)
	if err != nil {
		return mailproviderservices.SendMessageInput{}, coreerrors.NewError(mailprovidererror.InvalidMessageIdempotencyKey)
	}
	gid, err := validatorhelpers.NormalizeGID(input.GID)
	if err != nil {
		return mailproviderservices.SendMessageInput{}, err
	}
	senderId, err := validatorhelpers.NormalizeSenderId(input.SenderId)
	if err != nil {
		return mailproviderservices.SendMessageInput{}, err
	}
	to, err := validatorhelpers.NormalizeRecipients(input.To)
	if err != nil {
		return mailproviderservices.SendMessageInput{}, err
	}
	cc, err := validatorhelpers.NormalizeRecipients(input.Cc)
	if err != nil {
		return mailproviderservices.SendMessageInput{}, err
	}
	bcc, err := validatorhelpers.NormalizeRecipients(input.Bcc)
	if err != nil {
		return mailproviderservices.SendMessageInput{}, err
	}
	if err := validatorhelpers.ValidateMessageContent(to, cc, bcc, input.Subject, input.TextBody, input.HTMLBody); err != nil {
		return mailproviderservices.SendMessageInput{}, err
	}

	input.IdempotencyKey = idempotencyKey
	input.GID = string(gid)
	input.SenderId = string(senderId)
	input.To = to
	input.Cc = cc
	input.Bcc = bcc
	input.Subject = strings.TrimSpace(input.Subject)
	return input, nil
}

func (validator *MessageValidator) ValidateMessageId(value string) (entity.MessageId, error) {
	return validatorhelpers.NormalizeMessageId(value)
}

func (validator *MessageValidator) ValidateSearch(
	request coresearch.Request,
) (coresearch.Query[entity.MessageId], coresearch.CursorPayload, error) {
	return validatorhelpers.NormalizeMessageSearchRequest(request)
}
