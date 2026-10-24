package validators

import (
	"grovehub/services/mailProvider/entity"
	mailproviderservices "grovehub/services/mailProvider/services"
)

type DomainValidator interface {
	ValidateRegister(input mailproviderservices.RegisterDomainInput) (mailproviderservices.RegisterDomainInput, error)
	ValidateDomainId(value string) (entity.DomainId, error)
	ValidateDomainName(value string) (string, error)
}
