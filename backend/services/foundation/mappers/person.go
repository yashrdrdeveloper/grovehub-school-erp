package mappers

import (
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	"grovehub/services/foundation/entity"
)

type PersonMapper interface {
	ToPersonGID(request foundationrequest.PersonByGID) entity.GID
	ToPersonResponse(person entity.Person) foundationresponse.Person
	ToPersonErrorMetadata(gid entity.GID) foundationresponse.PersonErrorMetadata
}
