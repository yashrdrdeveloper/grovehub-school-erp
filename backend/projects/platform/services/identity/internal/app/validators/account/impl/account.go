package impl

import (
	"errors"
	"strings"

	coresearch "grovehub/libs/core/search"
	coreutils "grovehub/libs/core/utils"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	accountvalidator "grovehub/services/identity/internal/app/validators/account"
	searchhelpers "grovehub/services/identity/internal/app/validators/account/helpers"
)

// @Validator
type Validator struct {
	normalizedSearch bool
}

var _ accountvalidator.Validator = Validator{}

func New() accountvalidator.Validator {
	return Validator{}
}

func (Validator) NormalizeCreate(request accountdto.Create) (accountdto.Create, error) {
	email, err := coreutils.NormalizeEmail(request.Email)
	if err != nil {
		return request, errors.New("invalid email")
	}
	request.Email = email

	if strings.TrimSpace(request.PhoneNumber) != "" {
		phone, phoneErr := coreutils.NormalizePhone(request.PhoneNumber)
		if phoneErr != nil {
			return request, errors.New("invalid phone")
		}
		request.PhoneNumber = phone
	}

	request.Person.FirstName = strings.TrimSpace(request.Person.FirstName)
	request.Person.LastName = strings.TrimSpace(request.Person.LastName)
	if request.Person.MiddleName != nil {
		middleName := strings.TrimSpace(*request.Person.MiddleName)
		request.Person.MiddleName = &middleName
	}
	if request.Person.FirstName == "" || request.Person.LastName == "" ||
		len(request.Person.FirstName) > 100 || len(request.Person.LastName) > 100 {
		return request, errors.New("invalid registration")
	}
	if request.RoleAssignment != nil {
		assignment, err := searchhelpers.NormalizeRoleAssignment(*request.RoleAssignment)
		if err != nil {
			return request, err
		}
		request.RoleAssignment = &assignment
	}
	return request, nil
}

func (Validator) NormalizeUsername(username string) (string, error) {
	username = strings.ToLower(strings.TrimSpace(username))
	if len(username) < 3 || len(username) > 40 {
		return "", errors.New("invalid username")
	}
	for _, char := range username {
		letter := char >= 'a' && char <= 'z'
		digit := char >= '0' && char <= '9'
		punctuation := char == '.' || char == '_' || char == '-'
		if !letter && !digit && !punctuation {
			return "", errors.New("invalid username")
		}
	}
	return username, nil
}

func NewNormalized() accountvalidator.Validator {
	return Validator{normalizedSearch: true}
}

func (validator Validator) Search(request coresearch.Request) (coresearch.Query[string], error) {
	if validator.normalizedSearch {
		return searchhelpers.NormalizeNormalizedSearch(request)
	}
	return searchhelpers.NormalizeSearch(request)
}
