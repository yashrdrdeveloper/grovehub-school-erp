package helpers

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"

	coreerrors "grovehub/libs/core/errors"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	identityerror "grovehub/services/identity/constants/error"
	identityrequest "grovehub/services/identity/dto/request"
	"grovehub/services/identity/entity"
)

const (
	minPasswordLength = 12
	maxPasswordLength = 128
)

var phonePattern = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)

func ValidateBootstrapAccount(
	request identityrequest.BootstrapAccount,
) (entity.LoginIdentifierType, string, error) {
	if err := ValidatePassword(request.Password); err != nil {
		return "", "", err
	}

	identifierType, normalizedValue, err := NormalizeLoginIdentifier(request.Identifier)
	if err != nil {
		return "", "", err
	}
	if !PersonRequestContainsIdentifier(request.Person, identifierType, normalizedValue) {
		return "", "", invalidBootstrapRequest()
	}

	return identifierType, normalizedValue, nil
}

func ValidatePassword(password string) error {
	if !utf8.ValidString(password) {
		return invalidBootstrapRequest()
	}
	length := utf8.RuneCountInString(password)
	if length < minPasswordLength || length > maxPasswordLength {
		return invalidBootstrapRequest()
	}
	return nil
}

func NormalizeLoginIdentifier(
	identifier identityrequest.LoginIdentifier,
) (entity.LoginIdentifierType, string, error) {
	identifierType := entity.LoginIdentifierType(strings.ToUpper(strings.TrimSpace(identifier.Type)))
	normalizers := map[entity.LoginIdentifierType]func(string) (string, error){
		entity.LoginIdentifierTypeEmail: normalizeEmail,
		entity.LoginIdentifierTypePhone: normalizePhone,
	}
	normalize, ok := normalizers[identifierType]
	if !ok {
		return "", "", invalidBootstrapRequest()
	}

	normalizedValue, err := normalize(identifier.Value)
	if err != nil {
		return "", "", err
	}
	return identifierType, normalizedValue, nil
}

func PersonRequestContainsIdentifier(
	person foundationrequest.CreatePerson,
	identifierType entity.LoginIdentifierType,
	normalizedValue string,
) bool {
	groups := map[entity.LoginIdentifierType][]foundationrequest.ContactValue{
		entity.LoginIdentifierTypeEmail: person.Contacts.Emails,
		entity.LoginIdentifierTypePhone: person.Contacts.Phones,
	}

	for _, contact := range groups[identifierType] {
		candidate, err := normalizeContactValue(identifierType, contact.Value)
		if err == nil && candidate == normalizedValue {
			return true
		}
	}
	return false
}

func PersonIdentifierVerified(
	person foundationresponse.Person,
	identifierType entity.LoginIdentifierType,
	normalizedValue string,
) bool {
	groups := map[entity.LoginIdentifierType][]foundationresponse.ContactValue{
		entity.LoginIdentifierTypeEmail: person.Contacts.Emails,
		entity.LoginIdentifierTypePhone: person.Contacts.Phones,
	}

	for _, contact := range groups[identifierType] {
		candidate, err := normalizeContactValue(identifierType, contact.Value)
		if err == nil && candidate == normalizedValue {
			return contact.IsVerified
		}
	}
	return false
}

func normalizeContactValue(identifierType entity.LoginIdentifierType, value string) (string, error) {
	normalizers := map[entity.LoginIdentifierType]func(string) (string, error){
		entity.LoginIdentifierTypeEmail: normalizeEmail,
		entity.LoginIdentifierTypePhone: normalizePhone,
	}
	normalize, ok := normalizers[identifierType]
	if !ok {
		return "", invalidBootstrapRequest()
	}
	return normalize(value)
}

func normalizeEmail(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	address, err := mail.ParseAddress(normalized)
	if err != nil || strings.ToLower(address.Address) != normalized {
		return "", invalidBootstrapRequest()
	}
	return normalized, nil
}

func normalizePhone(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	builder := strings.Builder{}
	for index, character := range trimmed {
		digit := character >= '0' && character <= '9'
		separator := character == ' ' || character == '-' || character == '(' || character == ')'
		if character == '+' && index == 0 {
			builder.WriteRune(character)
		} else if digit {
			builder.WriteRune(character)
		} else if !separator {
			return "", invalidBootstrapRequest()
		}
	}

	normalized := builder.String()
	if !phonePattern.MatchString(normalized) {
		return "", invalidBootstrapRequest()
	}
	return normalized, nil
}

func invalidBootstrapRequest() error {
	return coreerrors.NewError(identityerror.InvalidBootstrapRequest)
}
