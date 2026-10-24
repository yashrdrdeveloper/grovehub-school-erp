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
) (identityvalidators.ValidatedCreateAccount, error) {
	personGid, err := validatorhelpers.NormalizePersonGid(input.PersonGid)
	if err != nil {
		return identityvalidators.ValidatedCreateAccount{}, err
	}
	email, err := validatorhelpers.NormalizeEmail(input.Email)
	if err != nil {
		return identityvalidators.ValidatedCreateAccount{}, err
	}
	phone, err := validatorhelpers.NormalizePhone(input.Phone)
	if err != nil {
		return identityvalidators.ValidatedCreateAccount{}, err
	}

	return identityvalidators.ValidatedCreateAccount{
		PersonGid: personGid,
		Email:     email,
		Phone:     phone,
	}, nil
}

func (validator *AccountValidator) ValidateUpdate(
	input identityservices.UpdateAccountInput,
) (identityvalidators.ValidatedUpdateAccount, error) {
	validated := identityvalidators.ValidatedUpdateAccount{}

	if input.Email != nil {
		email, err := validatorhelpers.NormalizeEmail(*input.Email)
		if err != nil {
			return identityvalidators.ValidatedUpdateAccount{}, err
		}
		validated.Email = &email
	}
	if input.Phone != nil {
		phone, err := validatorhelpers.NormalizePhone(*input.Phone)
		if err != nil {
			return identityvalidators.ValidatedUpdateAccount{}, err
		}
		validated.Phone = &phone
	}

	return validated, nil
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
