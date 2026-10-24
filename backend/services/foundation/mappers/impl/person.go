package impl

import (
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	"grovehub/services/foundation/entity"
	foundationmappers "grovehub/services/foundation/mappers"
)

type PersonMapper struct{}

func NewPersonMapper() foundationmappers.PersonMapper {
	return &PersonMapper{}
}

func (mapper *PersonMapper) ToPersonGID(request foundationrequest.PersonByGID) entity.GID {
	return entity.GID(request.GID)
}

func (mapper *PersonMapper) ToPersonResponse(person entity.Person) foundationresponse.Person {
	contacts := foundationresponse.Contacts{
		Emails: make([]foundationresponse.ContactValue, 0),
		Phones: make([]foundationresponse.ContactValue, 0),
	}

	appendContact := map[entity.ContactType]func(foundationresponse.ContactValue){
		entity.ContactTypeEmail: func(value foundationresponse.ContactValue) {
			contacts.Emails = append(contacts.Emails, value)
		},
		entity.ContactTypePhone: func(value foundationresponse.ContactValue) {
			contacts.Phones = append(contacts.Phones, value)
		},
	}

	for _, contact := range person.GetContacts() {
		appendValue, ok := appendContact[contact.GetType()]
		if !ok {
			continue
		}
		appendValue(foundationresponse.ContactValue{
			Value:      contact.GetValue(),
			IsPrimary:  contact.GetIsPrimary(),
			IsVerified: contact.GetIsVerified(),
		})
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

func (mapper *PersonMapper) ToPersonErrorMetadata(gid entity.GID) foundationresponse.PersonErrorMetadata {
	return foundationresponse.PersonErrorMetadata{GID: string(gid)}
}
