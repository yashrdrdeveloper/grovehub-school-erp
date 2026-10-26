package impl

import "grovehub/services/foundation/entity"

func toPersonModel(person entity.Person) personModel {
	var middleName *string
	if person.GetMiddleName() != "" {
		value := person.GetMiddleName()
		middleName = &value
	}

	return personModel{
		GID:        string(person.GetGID()),
		FirstName:  person.GetFirstName(),
		MiddleName: middleName,
		LastName:   person.GetLastName(),
		Status:     string(person.GetStatus()),
		Version:    person.GetVersion(),
		CreatedAt:  person.GetCreatedAt(),
		UpdatedAt:  person.GetUpdatedAt(),
	}
}

func toPersonContactModels(person entity.Person) []personContactModel {
	contacts := person.GetContacts()
	models := make([]personContactModel, 0, len(contacts))
	for _, contact := range contacts {
		models = append(models, personContactModel{
			PersonGID:       string(person.GetGID()),
			Type:            string(contact.GetType()),
			Value:           contact.GetValue(),
			NormalizedValue: contact.GetNormalizedValue(),
			IsPrimary:       contact.GetIsPrimary(),
			IsVerified:      contact.GetIsVerified(),
			CreatedAt:       person.GetUpdatedAt(),
			UpdatedAt:       person.GetUpdatedAt(),
		})
	}
	return models
}

func toPersonEntity(model personModel, contacts []personContactModel) entity.Person {
	middleName := ""
	if model.MiddleName != nil {
		middleName = *model.MiddleName
	}

	contactEntities := make([]entity.Contact, 0, len(contacts))
	for _, contact := range contacts {
		contactEntities = append(contactEntities, entity.NewContact(
			entity.ContactType(contact.Type),
			contact.Value,
			contact.NormalizedValue,
			contact.IsPrimary,
			contact.IsVerified,
		))
	}

	return entity.NewPerson(
		entity.GID(model.GID),
		model.FirstName,
		middleName,
		model.LastName,
		contactEntities,
		entity.PersonStatus(model.Status),
		model.Version,
		model.CreatedAt,
		model.UpdatedAt,
	)
}
