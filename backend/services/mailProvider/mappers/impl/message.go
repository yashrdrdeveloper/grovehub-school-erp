package impl

import (
	mailproviderresponse "grovehub/services/mailProvider/dto/response"
	"grovehub/services/mailProvider/entity"
	mailprovidermappers "grovehub/services/mailProvider/mappers"
)

type MessageMapper struct{}

func NewMessageMapper() mailprovidermappers.MessageMapper { return &MessageMapper{} }

func (mapper *MessageMapper) ToResponse(message entity.Message) mailproviderresponse.Message {
	return mailproviderresponse.Message{
		MessageId:           string(message.GetMessageId()),
		GID:                 string(message.GetGID()),
		SenderId:            string(message.GetSenderId()),
		ThreadId:            string(message.GetThreadId()),
		ParentMessageId:     string(message.GetParentMessageId()),
		ResendOfMessageId:   string(message.GetResendOfMessageId()),
		InternetMessageId:   message.GetInternetMessageId(),
		InReplyTo:           message.GetInReplyTo(),
		ReferenceMessageIds: message.GetReferenceMessageIds(),
		To:                  message.GetTo(),
		Cc:                  message.GetCc(),
		Bcc:                 message.GetBcc(),
		Subject:             message.GetSubject(),
		Status:              string(message.GetStatus()),
		FailureMessage:      message.GetFailureMessage(),
		Version:             message.GetVersion(),
		CreatedAt:           message.GetCreatedAt(),
		UpdatedAt:           message.GetUpdatedAt(),
	}
}
