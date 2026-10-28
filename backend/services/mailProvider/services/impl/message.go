package impl

import (
	"context"
	"errors"
	"time"

	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/id"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
	"grovehub/services/mailProvider/entity"
	mailproviderrepository "grovehub/services/mailProvider/repository"
	mailproviderservices "grovehub/services/mailProvider/services"
	mailprovidertransport "grovehub/services/mailProvider/transport"
	mailprovidervalidators "grovehub/services/mailProvider/validators"
)

type MessageService struct {
	repository    mailproviderrepository.MessageRepository
	domainService mailproviderservices.DomainService
	senderService mailproviderservices.SenderService
	transport     mailprovidertransport.Sender
	validator     mailprovidervalidators.MessageValidator
}

func NewMessageService(
	repository mailproviderrepository.MessageRepository,
	domainService mailproviderservices.DomainService,
	senderService mailproviderservices.SenderService,
	transport mailprovidertransport.Sender,
	validator mailprovidervalidators.MessageValidator,
) mailproviderservices.MessageService {
	return &MessageService{repository: repository, domainService: domainService, senderService: senderService, transport: transport, validator: validator}
}

func (service *MessageService) Send(ctx context.Context, input mailproviderservices.SendMessageInput) (entity.Message, error) {
	validated, err := service.validator.ValidateSend(input)
	if err != nil {
		return entity.Message{}, err
	}
	sender, err := service.senderService.GetById(ctx, validated.SenderId)
	if err != nil {
		return entity.Message{}, err
	}
	if sender.GetStatus() != entity.SenderStatusActive {
		return entity.Message{}, coreerrors.NewError(mailprovidererror.InvalidSender)
	}
	domain, err := service.domainService.GetById(ctx, string(sender.GetDomainId()))
	if err != nil {
		return entity.Message{}, err
	}
	if domain.GetStatus() != entity.DomainStatusVerified {
		return entity.Message{}, coreerrors.NewError(mailprovidererror.DomainNotVerified)
	}
	messageId, err := coreid.NewUUIDv7()
	if err != nil {
		return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageCreateFailed, err)
	}
	now := time.Now().UTC()
	message := entity.NewMessage(entity.MessageId(messageId), sender.GetSenderId(), validated.To, validated.Cc, validated.Bcc,
		validated.Subject, validated.TextBody, validated.HTMLBody, entity.MessageStatusQueued, "", 1, now, now)
	if err := service.repository.Create(ctx, message); err != nil {
		return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageCreateFailed, err)
	}
	if err := service.transport.Send(ctx, message, sender); err != nil {
		failed := message.SetStatus(entity.MessageStatusFailed).SetFailureMessage(err.Error()).SetVersion(2).SetUpdatedAt(time.Now().UTC())
		if updateErr := service.repository.Update(ctx, failed); updateErr != nil {
			return failed, coreerrors.Wrap(mailprovidererror.MessageUpdateFailed, errors.Join(err, updateErr))
		}
		return failed, coreerrors.Wrap(mailprovidererror.MessageSubmissionFailed, err)
	}
	submitted := message.SetStatus(entity.MessageStatusSubmitted).SetVersion(2).SetUpdatedAt(time.Now().UTC())
	if err := service.repository.Update(ctx, submitted); err != nil {
		return submitted, coreerrors.Wrap(mailprovidererror.MessageUpdateFailed, err)
	}
	return submitted, nil
}

func (service *MessageService) GetById(ctx context.Context, messageId string) (entity.Message, error) {
	validatedId, err := service.validator.ValidateMessageId(messageId)
	if err != nil {
		return entity.Message{}, err
	}
	message, err := service.repository.FindById(ctx, validatedId)
	if errors.Is(err, mailproviderrepository.ErrMessageNotFound) {
		return entity.Message{}, coreerrors.NewError(mailprovidererror.MessageNotFound)
	}
	if err != nil {
		return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageReadFailed, err)
	}
	return message, nil
}
