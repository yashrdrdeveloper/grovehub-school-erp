package mappers

import (
	foundationresponse "grovehub/services/foundation/dto/response"
	"grovehub/services/foundation/entity"
)

func ToPersonErrorMetadata(gid entity.GID) foundationresponse.PersonErrorMetadata {
	return foundationresponse.PersonErrorMetadata{
		GID: string(gid),
	}
}
