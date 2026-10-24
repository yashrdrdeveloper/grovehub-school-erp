package mappers

import (
	mailproviderresponse "grovehub/services/mailProvider/dto/response"
	"grovehub/services/mailProvider/entity"
)

type SenderMapper interface {
	ToResponse(sender entity.Sender) mailproviderresponse.Sender
}
