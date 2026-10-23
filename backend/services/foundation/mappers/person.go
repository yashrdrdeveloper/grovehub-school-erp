package mappers

import (
	"time"

	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	"grovehub/services/foundation/entity"
)

type PersonBuilder struct {
	person entity.Person
}

func NewPersonBuilder() *PersonBuilder {
	return &PersonBuilder{}
}

func NewPersonBuilderFrom(person entity.Person) *PersonBuilder {
	return &PersonBuilder{person: person}
}

func (builder *PersonBuilder) WithGID(gid entity.GID) *PersonBuilder {
	builder.person.GID = gid
	return builder
}

func (builder *PersonBuilder) WithStatus(status entity.PersonStatus) *PersonBuilder {
	builder.person.Status = status
	return builder
}

func (builder *PersonBuilder) WithVersion(version int64) *PersonBuilder {
	builder.person.Version = version
	return builder
}

func (builder *PersonBuilder) WithCreatedAt(createdAt time.Time) *PersonBuilder {
	builder.person.CreatedAt = createdAt
	return builder
}

func (builder *PersonBuilder) WithUpdatedAt(updatedAt time.Time) *PersonBuilder {
	builder.person.UpdatedAt = updatedAt
	return builder
}

func (builder *PersonBuilder) Build() entity.Person {
	return builder.person
}

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
