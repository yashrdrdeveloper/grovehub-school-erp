package impl

import (
	"grovehub/services/mailProvider/entity"
	mailproviderservices "grovehub/services/mailProvider/services"
	mailprovidervalidators "grovehub/services/mailProvider/validators"
	validatorhelpers "grovehub/services/mailProvider/validators/helpers"
)

type DomainValidator struct{}

func NewDomainValidator() mailprovidervalidators.DomainValidator {
	return &DomainValidator{}
}

func (validator *DomainValidator) ValidateRegister(
	input mailproviderservices.RegisterDomainInput,
) (mailproviderservices.RegisterDomainInput, error) {
	name, err := validatorhelpers.NormalizeDomain(input.Name)
	if err != nil {
		return mailproviderservices.RegisterDomainInput{}, err
	}
	input.Name = name
	return input, nil
}

func (validator *DomainValidator) ValidateDomainId(value string) (entity.DomainId, error) {
	return validatorhelpers.NormalizeDomainId(value)
}

func (validator *DomainValidator) ValidateDomainName(value string) (string, error) {
	return validatorhelpers.NormalizeDomain(value)
}
