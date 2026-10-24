package impl

import (
	"strings"

	"grovehub/services/mailProvider/entity"
	mailproviderservices "grovehub/services/mailProvider/services"
	mailprovidervalidators "grovehub/services/mailProvider/validators"
	validatorhelpers "grovehub/services/mailProvider/validators/helpers"
)

type MessageValidator struct{}

func NewMessageValidator() mailprovidervalidators.MessageValidator {
	return &MessageValidator{}
}

func (validator *MessageValidator) ValidateSend(
	input mailproviderservices.SendMessageInput,
) (mailproviderservices.SendMessageInput, error) {
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
