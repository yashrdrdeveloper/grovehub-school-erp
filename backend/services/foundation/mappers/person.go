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
	return foundationresponse.Person{
		GID:       string(person.GID),
		Status:    string(person.Status),
		Version:   person.Version,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
	}
}
