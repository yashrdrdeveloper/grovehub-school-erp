package helpers

import (
    coreerrors "grovehub/libs/core/errors"
    corenormalize "grovehub/libs/core/normalize"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    persondto "grovehub/services/foundation/client/person/dto"
    "grovehub/services/foundation/internal/app/entity"
)

const maxContactsPerType = 10

// NormalizeContacts canonicalizes contact values only. Cross-field rules
// such as duplicate and primary-contact limits run after DTO tag validation.
func NormalizeContacts(request persondto.Contacts) (persondto.Contacts, error) {
	emails, err := normalizeContactGroup(request.Emails, normalizeEmail)
	if err != nil {
		return persondto.Contacts{}, err
	}
	phones, err := normalizeContactGroup(request.Phones, normalizePhone)
	if err != nil {
		return persondto.Contacts{}, err
	}
	return persondto.Contacts{Emails: emails, Phones: phones}, nil
}

// ValidateContacts enforces request-level contact rules on normalized values.
func ValidateContacts(request persondto.Contacts) error {
	if len(request.Emails) > maxContactsPerType || len(request.Phones) > maxContactsPerType {
		return coreerrors.NewError(personerror.InvalidPersonContact)
	}
	seen := make(map[string]struct{}, len(request.Emails)+len(request.Phones))
	for _, group := range []struct {
		kind entity.ContactType
		values []persondto.ContactValue
	}{
		{entity.ContactTypeEmail, request.Emails},
		{entity.ContactTypePhone, request.Phones},
	} {
		primaryCount := 0
		for _, contact := range group.values {
			key := contactKey(group.kind, contact.Value)
			if _, duplicate := seen[key]; duplicate {
				return coreerrors.NewError(personerror.InvalidPersonContact)
			}
			seen[key] = struct{}{}
			if contact.IsPrimary {
				primaryCount++
				if primaryCount > 1 {
					return coreerrors.NewError(personerror.InvalidPersonContact)
				}
			}
		}
	}
	return nil
}

func normalizeContactGroup(
	values []persondto.ContactValue,
	normalize func(string) (string, error),
) ([]persondto.ContactValue, error) {
	normalized := make([]persondto.ContactValue, 0, len(values))
	for _, contact := range values {
		value, err := normalize(contact.Value)
		if err != nil {
			return nil, err
		}
		normalized = append(normalized, persondto.ContactValue{
			Value: value,
			IsPrimary: contact.IsPrimary,
		})
	}
	return normalized, nil
}

func normalizeEmail(value string) (string, error) {
	normalized, err := corenormalize.Email(value)
	if err != nil {
		return "", coreerrors.NewError(personerror.InvalidPersonContact)
	}
	return normalized, nil
}

func normalizePhone(value string) (string, error) {
	normalized, err := corenormalize.Phone(value)
	if err != nil {
		return "", coreerrors.NewError(personerror.InvalidPersonContact)
	}
	return normalized, nil
}

func contactKey(contactType entity.ContactType, normalizedValue string) string {
	return string(contactType) + "\x00" + normalizedValue
}
