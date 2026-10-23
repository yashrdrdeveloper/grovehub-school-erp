package mappers

import (
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	"grovehub/services/foundation/entity"
)

func ToPersonGID(request foundationrequest.PersonByGID) entity.GID {
	return entity.GID(request.GID)
}

func ToPersonResponse(person entity.Person) foundationresponse.Person {
	contacts := foundationresponse.Contacts{
		Emails: make([]foundationresponse.ContactValue, 0),
		Phones: make([]foundationresponse.ContactValue, 0),
	}

	for _, contact := range person.GetContacts() {
		value := foundationresponse.ContactValue{
			Value:      contact.GetValue(),
			IsPrimary:  contact.GetIsPrimary(),
			IsVerified: contact.GetIsVerified(),
		}

		switch contact.GetType() {
		case entity.ContactTypeEmail:
			contacts.Emails = append(contacts.Emails, value)
		case entity.ContactTypePhone:
			contacts.Phones = append(contacts.Phones, value)
		}
	}

	return foundationresponse.Person{
		GID:        string(person.GetGID()),
		FirstName:  person.GetFirstName(),
		MiddleName: person.GetMiddleName(),
		LastName:   person.GetLastName(),
		Contacts:   contacts,
		Status:     string(person.GetStatus()),
		Version:    person.GetVersion(),
		CreatedAt:  person.GetCreatedAt(),
		UpdatedAt:  person.GetUpdatedAt(),
	}
}
