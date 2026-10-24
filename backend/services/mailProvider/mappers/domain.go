package mappers

import (
	mailproviderresponse "grovehub/services/mailProvider/dto/response"
	"grovehub/services/mailProvider/entity"
)

type DomainMapper interface {
	ToResponse(domain entity.Domain) mailproviderresponse.Domain
}
