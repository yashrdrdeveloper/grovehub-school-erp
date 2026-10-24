package mappers

import (
	mailproviderresponse "grovehub/services/mailProvider/dto/response"
	"grovehub/services/mailProvider/entity"
)

type MessageMapper interface {
	ToResponse(message entity.Message) mailproviderresponse.Message
}
