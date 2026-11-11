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

func (mapper *PersonMapper) ToPersonGID(request foundationrequest.PersonByGID) string {
	return request.GID
}

func (mapper *PersonMapper) ToPersonResponse(person entity.Person, personContacts []entity.Contact) foundationresponse.Person {
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

	for _, contact := range personContacts {
		appendValue, ok := appendContact[contact.Type]
		if !ok {
			continue
		}
		appendValue(foundationresponse.ContactValue{
			Value:      contact.Value,
			IsPrimary:  contact.IsPrimary,
			IsVerified: contact.IsVerified,
		})
	}

	middleName := ""
	if person.MiddleName != nil { middleName = *person.MiddleName }
	return foundationresponse.Person{
		GID:        person.GID,
		FirstName:  person.FirstName,
		MiddleName: middleName,
		LastName:   person.LastName,
		Contacts:   contacts,
		Status:     string(person.Status),
		Version:    person.Version,
		CreatedAt:  person.CreatedAt,
		UpdatedAt:  person.UpdatedAt,
	}
}

