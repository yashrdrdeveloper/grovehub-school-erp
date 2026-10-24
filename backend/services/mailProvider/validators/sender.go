package validators

import (
	"grovehub/services/mailProvider/entity"
	mailproviderservices "grovehub/services/mailProvider/services"
)

type SenderValidator interface {
	ValidateRegister(input mailproviderservices.RegisterSenderInput) (mailproviderservices.RegisterSenderInput, error)
	ValidateSenderId(value string) (entity.SenderId, error)
}
