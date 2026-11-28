package person

import (
	persondto "grovehub/services/foundation/internal/app/dto/person"
	"grovehub/services/foundation/internal/app/entity"
)

// @Mapper PersonMapper
type PersonMapper interface {
	ToPersonGID(request persondto.PersonByGID) string
	ToPersonResponse(person entity.Person, contacts []entity.Contact) persondto.Person
	// @Mapping GID=gid
	ToPersonErrorMetadata(gid string) persondto.PersonErrorMetadata
}
