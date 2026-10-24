package impl

import (
	mailproviderresponse "grovehub/services/mailProvider/dto/response"
	"grovehub/services/mailProvider/entity"
	mailprovidermappers "grovehub/services/mailProvider/mappers"
	servicehelpers "grovehub/services/mailProvider/services/helpers"
)

type DomainMapper struct{}

func NewDomainMapper() mailprovidermappers.DomainMapper { return &DomainMapper{} }

func (mapper *DomainMapper) ToResponse(domain entity.Domain) mailproviderresponse.Domain {
	return mailproviderresponse.Domain{
		DomainId:          string(domain.GetDomainId()),
		Name:              domain.GetName(),
		VerificationHost:  servicehelpers.DomainVerificationHost(domain.GetName()),
		VerificationValue: servicehelpers.DomainVerificationValue(domain.GetVerificationToken()),
		Status:            string(domain.GetStatus()),
		Version:           domain.GetVersion(),
		CreatedAt:         domain.GetCreatedAt(),
		UpdatedAt:         domain.GetUpdatedAt(),
	}
}
