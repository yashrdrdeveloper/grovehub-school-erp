package impl

import (
	"context"
	"errors"
	"time"

	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/id"
	coreidempotency "grovehub/libs/core/idempotency"
	coresearch "grovehub/libs/core/search"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
	"grovehub/services/mailProvider/entity"
	mailproviderrepository "grovehub/services/mailProvider/repository"
	mailproviderservices "grovehub/services/mailProvider/services"
	servicehelpers "grovehub/services/mailProvider/services/helpers"
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
	return &MessageService{
		repository: repository, domainService: domainService, senderService: senderService,
		transport: transport, validator: validator,
	}
}

func (service *MessageService) Send(
	ctx context.Context,
	input mailproviderservices.SendMessageInput,
) (entity.Message, error) {
	validated, err := service.validator.ValidateSend(input)
	if err != nil {
		return entity.Message{}, err
	}
	requestHash, err := coreidempotency.Fingerprint(struct {
		GID      string   `json:"gid"`
		SenderId string   `json:"senderId"`
		To       []string `json:"to"`
		Cc       []string `json:"cc"`
		Bcc      []string `json:"bcc"`
		Subject  string   `json:"subject"`
		TextBody string   `json:"textBody"`
		HTMLBody string   `json:"htmlBody"`
	}{
		GID: validated.GID, SenderId: validated.SenderId,
		To: validated.To, Cc: validated.Cc, Bcc: validated.Bcc,
		Subject: validated.Subject, TextBody: validated.TextBody, HTMLBody: validated.HTMLBody,
	})
	if err != nil {
		return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageCreateFailed, err)
	}

	existing, err := service.repository.FindByIdempotencyKey(ctx, validated.IdempotencyKey)
	if err == nil {
		return service.replayOrSubmit(ctx, existing, requestHash)
	}
	if !errors.Is(err, mailproviderrepository.ErrMessageNotFound) {
		return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageReadFailed, err)
	}

	messageId, err := coreid.NewUUIDv7()
	if err != nil {
		return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageCreateFailed, err)
	}
	now := time.Now().UTC()
	message := entity.NewMessage(
		entity.MessageId(messageId),
		validated.IdempotencyKey,
		requestHash,
		entity.GID(validated.GID),
		entity.SenderId(validated.SenderId),
		validated.To,
		validated.Cc,
		validated.Bcc,
		validated.Subject,
		validated.TextBody,
		validated.HTMLBody,
		entity.MessageStatusQueued,
		"",
		1,
		now,
		now,
	)
	if err := service.repository.Create(ctx, message); err != nil {
		if errors.Is(err, mailproviderrepository.ErrMessageIdempotencyKeyAlreadyExists) {
			existing, readErr := service.repository.FindByIdempotencyKey(ctx, validated.IdempotencyKey)
			if readErr != nil {
				return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageReadFailed, readErr)
			}
			return service.replayOrSubmit(ctx, existing, requestHash)
		}
		return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageCreateFailed, err)
	}
	return service.submit(ctx, message)
}

func (service *MessageService) Retry(ctx context.Context, messageId string) (entity.Message, error) {
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
	if message.GetStatus() != entity.MessageStatusFailed {
		return entity.Message{}, coreerrors.NewError(mailprovidererror.MessageRetryNotAllowed)
	}

	sender, err := service.resolveSubmissionSender(ctx, message)
	if err != nil {
		return entity.Message{}, err
	}

	prepared, err := service.repository.PrepareRetry(ctx, message, time.Now().UTC())
	if err != nil {
		if errors.Is(err, mailproviderrepository.ErrMessageVersionConflict) {
			latest, readErr := service.repository.FindById(ctx, validatedId)
			if readErr != nil {
				return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageReadFailed, readErr)
			}
			if latest.GetStatus() == entity.MessageStatusSubmitted {
				return latest, nil
			}
			if latest.GetStatus() == entity.MessageStatusQueued && latest.GetSubmissionStartedAt() != nil {
				return latest, coreerrors.NewError(mailprovidererror.MessageSubmissionInProgress)
			}
			return entity.Message{}, coreerrors.NewError(mailprovidererror.MessageConcurrentModification)
		}
		return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageUpdateFailed, err)
	}

	return service.submitWithSender(ctx, prepared, sender)
}

func (service *MessageService) replayOrSubmit(
	ctx context.Context,
	message entity.Message,
	requestHash string,
) (entity.Message, error) {
	if message.GetRequestHash() != requestHash {
		return entity.Message{}, coreerrors.NewError(mailprovidererror.MessageIdempotencyConflict)
	}
	switch message.GetStatus() {
	case entity.MessageStatusSubmitted:
		return message, nil
	case entity.MessageStatusFailed:
		return message, coreerrors.NewError(mailprovidererror.MessageSubmissionFailed)
	case entity.MessageStatusQueued:
		if message.GetSubmissionStartedAt() != nil {
			return message, coreerrors.NewError(mailprovidererror.MessageSubmissionInProgress)
		}
		return service.submit(ctx, message)
	default:
		return entity.Message{}, coreerrors.NewError(mailprovidererror.InvalidMessage)
	}
}

func (service *MessageService) submit(ctx context.Context, message entity.Message) (entity.Message, error) {
	sender, err := service.resolveSubmissionSender(ctx, message)
	if err != nil {
		return entity.Message{}, err
	}
	return service.submitWithSender(ctx, message, sender)
}

func (service *MessageService) resolveSubmissionSender(
	ctx context.Context,
	message entity.Message,
) (entity.Sender, error) {
	sender, err := service.senderService.GetById(ctx, string(message.GetSenderId()))
	if err != nil {
		return entity.Sender{}, err
	}
	if sender.GetStatus() != entity.SenderStatusActive {
		return entity.Sender{}, coreerrors.NewError(mailprovidererror.InvalidSender)
	}
	domain, err := service.domainService.GetById(ctx, string(sender.GetDomainId()))
	if err != nil {
		return entity.Sender{}, err
	}
	if domain.GetStatus() != entity.DomainStatusVerified {
		return entity.Sender{}, coreerrors.NewError(mailprovidererror.DomainNotVerified)
	}
	return sender, nil
}

func (service *MessageService) submitWithSender(
	ctx context.Context,
	message entity.Message,
	sender entity.Sender,
) (entity.Message, error) {
	claimed, err := service.repository.ClaimSubmission(ctx, message, time.Now().UTC())
	if err != nil {
		if errors.Is(err, mailproviderrepository.ErrMessageSubmissionAlreadyClaimed) {
			latest, readErr := service.repository.FindById(ctx, message.GetMessageId())
			if readErr != nil {
				return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageReadFailed, readErr)
			}
			return service.replayOrSubmit(ctx, latest, message.GetRequestHash())
		}
		return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageUpdateFailed, err)
	}

	if err := service.transport.Send(ctx, claimed, sender); err != nil {
		failed := claimed.
			SetStatus(entity.MessageStatusFailed).
			SetFailureMessage(err.Error()).
			SetVersion(claimed.GetVersion() + 1).
			SetUpdatedAt(time.Now().UTC())
		if updateErr := service.repository.Update(ctx, failed); updateErr != nil {
			return failed, coreerrors.Wrap(mailprovidererror.MessageUpdateFailed, errors.Join(err, updateErr))
		}
		return failed, coreerrors.Wrap(mailprovidererror.MessageSubmissionFailed, err)
	}

	submitted := claimed.
		SetStatus(entity.MessageStatusSubmitted).
		SetVersion(claimed.GetVersion() + 1).
		SetUpdatedAt(time.Now().UTC())
	if err := service.repository.Update(ctx, submitted); err != nil {
		if errors.Is(err, mailproviderrepository.ErrMessageVersionConflict) {
			latest, readErr := service.repository.FindById(ctx, claimed.GetMessageId())
			if readErr == nil && latest.GetStatus() == entity.MessageStatusSubmitted {
				return latest, nil
			}
			if readErr != nil {
				return entity.Message{}, coreerrors.Wrap(mailprovidererror.MessageReadFailed, readErr)
			}
			return entity.Message{}, coreerrors.NewError(mailprovidererror.MessageConcurrentModification)
		}
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

func (service *MessageService) Search(
	ctx context.Context,
	request coresearch.Request,
) (coresearch.Result[entity.Message], error) {
	query, cursorPayload, err := service.validator.ValidateSearch(request)
	if err != nil {
		return coresearch.Result[entity.Message]{}, err
	}
	result, err := service.repository.Search(ctx, query)
	if err != nil {
		if errors.Is(err, coresearch.ErrInvalidQuery) {
			return coresearch.Result[entity.Message]{}, coreerrors.NewError(mailprovidererror.InvalidMessageSearch)
		}
		return coresearch.Result[entity.Message]{}, coreerrors.Wrap(mailprovidererror.MessageSearchFailed, err)
	}
	nextCursor, err := buildMessageNextCursor(result, query, cursorPayload)
	if err != nil {
		return coresearch.Result[entity.Message]{}, coreerrors.Wrap(mailprovidererror.MessageSearchFailed, err)
	}
	return coresearch.Result[entity.Message]{
		Items: result.Items,
		PageInfo: coresearch.PageResponse{
			Rows: query.Limit, NextCursor: nextCursor, HasNext: result.HasNext,
		},
	}, nil
}

func buildMessageNextCursor(
	result coresearch.QueryResult[entity.Message],
	query coresearch.Query[entity.MessageId],
	cursor coresearch.CursorPayload,
) (string, error) {
	if !result.HasNext || len(result.Items) == 0 {
		return "", nil
	}
	lastMessage := result.Items[len(result.Items)-1]
	cursor.Value = servicehelpers.MessageCursorValue(lastMessage, query.Sort.Field)
	cursor.CursorId = string(lastMessage.GetMessageId())
	return servicehelpers.EncodeMessageSearchCursor(cursor)
}
