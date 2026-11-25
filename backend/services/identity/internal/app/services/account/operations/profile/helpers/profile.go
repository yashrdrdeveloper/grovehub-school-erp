package helpers

import (
	"errors"
	"strings"
	"unicode/utf8"

	annotationdto "grovehub/libs/core/annotation/dto"
	foundation "grovehub/services/foundation/client/person"
	accountdto "grovehub/services/identity/internal/app/dto/account"
)

func MapPerson(person foundation.Person) accountdto.Profile {
	profile := accountdto.Profile{FirstName: person.FirstName, LastName: person.LastName}
	if person.MiddleName != "" {
		middle := person.MiddleName
		profile.MiddleName = &middle
	}
	for _, phone := range person.Contacts.Phones {
		if phone.IsPrimary {
			value := phone.Value
			profile.PhoneNumber = &value
			break
		}
	}
	return profile
}

func BuildPatch(input accountdto.UpdateProfile) (foundation.UpdatePerson, error) {
	first := strings.TrimSpace(input.FirstName)
	last := strings.TrimSpace(input.LastName)
	if first == "" || last == "" || utf8.RuneCountInString(first) > 100 || utf8.RuneCountInString(last) > 100 {
		return foundation.UpdatePerson{}, errors.New("invalid person name")
	}
	var middle *string
	if input.MiddleName != nil {
		value := strings.TrimSpace(*input.MiddleName)
		if utf8.RuneCountInString(value) > 100 {
			return foundation.UpdatePerson{}, errors.New("invalid middle name")
		}
		if value != "" {
			middle = &value
		}
	}
	return foundation.UpdatePerson{
		FirstName: annotationdto.Patch[string]{Set: true, Value: &first},
		MiddleName: annotationdto.Patch[string]{Set: true, Value: middle},
		LastName: annotationdto.Patch[string]{Set: true, Value: &last},
	}, nil
}
