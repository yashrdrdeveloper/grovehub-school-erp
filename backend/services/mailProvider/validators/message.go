package validators

import (
	"grovehub/services/mailProvider/entity"
	mailproviderservices "grovehub/services/mailProvider/services"
)

type MessageValidator interface {
	ValidateSend(input mailproviderservices.SendMessageInput) (mailproviderservices.SendMessageInput, error)
	ValidateMessageId(value string) (entity.MessageId, error)
}
