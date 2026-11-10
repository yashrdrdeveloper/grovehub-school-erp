package mappers

import (
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	"grovehub/services/foundation/entity"
)

// @Mapper PersonMapper
type PersonMapper interface {
	ToPersonGID(request foundationrequest.PersonByGID) string
	ToPersonResponse(person entity.Person) foundationresponse.Person
	// @Mapping GID=gid
	ToPersonErrorMetadata(gid string) foundationresponse.PersonErrorMetadata
}
