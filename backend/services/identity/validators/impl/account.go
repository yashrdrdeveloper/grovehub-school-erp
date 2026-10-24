package impl

import (
	"grovehub/services/identity/entity"
	identityservices "grovehub/services/identity/services"
	identityvalidators "grovehub/services/identity/validators"
	validatorhelpers "grovehub/services/identity/validators/helpers"
)

type AccountValidator struct{}

func NewAccountValidator() identityvalidators.AccountValidator {
	return &AccountValidator{}
}

func (validator *AccountValidator) ValidateCreate(
	input identityservices.CreateAccountInput,
) (identityservices.CreateAccountInput, error) {
	personGid, err := validatorhelpers.NormalizePersonGid(input.PersonGid)
	if err != nil {
		return identityservices.CreateAccountInput{}, err
	}
	email, err := validatorhelpers.NormalizeEmail(input.Email)
	if err != nil {
		return identityservices.CreateAccountInput{}, err
	}
	phone, err := validatorhelpers.NormalizePhone(input.Phone)
	if err != nil {
		return identityservices.CreateAccountInput{}, err
	}

	input.PersonGid = string(personGid)
	input.Email = email
	input.Phone = phone
	return input, nil
}

func (validator *AccountValidator) ValidateUpdate(
	input identityservices.UpdateAccountInput,
) (identityservices.UpdateAccountInput, error) {
	if input.Email != nil {
		email, err := validatorhelpers.NormalizeEmail(*input.Email)
		if err != nil {
			return identityservices.UpdateAccountInput{}, err
		}
		input.Email = &email
	}
	if input.Phone != nil {
		phone, err := validatorhelpers.NormalizePhone(*input.Phone)
		if err != nil {
			return identityservices.UpdateAccountInput{}, err
		}
		input.Phone = &phone
	}

	return input, nil
}

func (validator *AccountValidator) ValidateAccountId(value string) (entity.AccountId, error) {
	return validatorhelpers.NormalizeAccountId(value)
}

func (validator *AccountValidator) ValidatePersonGid(value string) (entity.PersonGid, error) {
	return validatorhelpers.NormalizePersonGid(value)
}

func (validator *AccountValidator) ValidateEmail(value string) (string, error) {
	return validatorhelpers.NormalizeEmail(value)
}
