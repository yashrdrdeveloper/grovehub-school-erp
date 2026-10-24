package impl

import (
	"grovehub/services/mailProvider/entity"
	mailproviderservices "grovehub/services/mailProvider/services"
	mailprovidervalidators "grovehub/services/mailProvider/validators"
	validatorhelpers "grovehub/services/mailProvider/validators/helpers"
)

type SenderValidator struct{}

func NewSenderValidator() mailprovidervalidators.SenderValidator {
	return &SenderValidator{}
}

func (validator *SenderValidator) ValidateRegister(
	input mailproviderservices.RegisterSenderInput,
) (mailproviderservices.RegisterSenderInput, error) {
	domainId, err := validatorhelpers.NormalizeDomainId(input.DomainId)
	if err != nil {
		return mailproviderservices.RegisterSenderInput{}, err
	}
	email, err := validatorhelpers.NormalizeSenderEmail(input.Email)
	if err != nil {
		return mailproviderservices.RegisterSenderInput{}, err
	}
	displayName, err := validatorhelpers.NormalizeDisplayName(input.DisplayName)
	if err != nil {
		return mailproviderservices.RegisterSenderInput{}, err
	}
	input.DomainId = string(domainId)
	input.Email = email
	input.DisplayName = displayName
	return input, nil
}

func (validator *SenderValidator) ValidateSenderId(value string) (entity.SenderId, error) {
	return validatorhelpers.NormalizeSenderId(value)
}
