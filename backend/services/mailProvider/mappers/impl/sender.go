package impl

import (
	mailproviderresponse "grovehub/services/mailProvider/dto/response"
	"grovehub/services/mailProvider/entity"
	mailprovidermappers "grovehub/services/mailProvider/mappers"
)

type SenderMapper struct{}

func NewSenderMapper() mailprovidermappers.SenderMapper { return &SenderMapper{} }

func (mapper *SenderMapper) ToResponse(sender entity.Sender) mailproviderresponse.Sender {
	return mailproviderresponse.Sender{
		SenderId:    string(sender.GetSenderId()),
		DomainId:    string(sender.GetDomainId()),
		Email:       sender.GetEmail(),
		DisplayName: sender.GetDisplayName(),
		Status:      string(sender.GetStatus()),
		Version:     sender.GetVersion(),
		CreatedAt:   sender.GetCreatedAt(),
		UpdatedAt:   sender.GetUpdatedAt(),
	}
}
