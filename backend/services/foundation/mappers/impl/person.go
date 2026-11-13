package impl

import (
	persondto "grovehub/services/foundation/dto/person"
	"grovehub/services/foundation/entity"
	foundationmappers "grovehub/services/foundation/mappers"
)

type PersonMapper struct{}

func NewPersonMapper() foundationmappers.PersonMapper {
	return &PersonMapper{}
}

func (mapper *PersonMapper) ToPersonGID(request persondto.PersonByGID) string {
	return request.GID
}

func (mapper *PersonMapper) ToPersonResponse(person entity.Person, personContacts []entity.Contact) persondto.Person {
	contacts := persondto.ResponseContacts{
		Emails: make([]persondto.ResponseContactValue, 0),
		Phones: make([]persondto.ResponseContactValue, 0),
	}

	appendContact := map[entity.ContactType]func(persondto.ResponseContactValue){
		entity.ContactTypeEmail: func(value persondto.ResponseContactValue) {
			contacts.Emails = append(contacts.Emails, value)
		},
		entity.ContactTypePhone: func(value persondto.ResponseContactValue) {
			contacts.Phones = append(contacts.Phones, value)
		},
	}

	for _, contact := range personContacts {
		appendValue, ok := appendContact[contact.Type]
		if !ok {
			continue
		}
		appendValue(persondto.ResponseContactValue{
			Value:      contact.Value,
			IsPrimary:  contact.IsPrimary,
			IsVerified: contact.IsVerified,
		})
	}

	middleName := ""
	if person.MiddleName != nil { middleName = *person.MiddleName }
	return persondto.Person{
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

