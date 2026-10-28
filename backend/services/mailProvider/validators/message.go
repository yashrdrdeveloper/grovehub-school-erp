package validators

import (
	coresearch "grovehub/libs/core/search"
	"grovehub/services/mailProvider/entity"
	mailproviderservices "grovehub/services/mailProvider/services"
)

type MessageValidator interface {
	ValidateSend(input mailproviderservices.SendMessageInput) (mailproviderservices.SendMessageInput, error)
	ValidateMessageId(value string) (entity.MessageId, error)
	ValidateSearch(request coresearch.Request) (coresearch.Query[entity.MessageId], coresearch.CursorPayload, error)
}
