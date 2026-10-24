package impl

import (
	"context"
	"errors"
	"strings"
	"time"

	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/helpers/id"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
	"grovehub/services/mailProvider/entity"
	mailproviderrepository "grovehub/services/mailProvider/repository"
	mailproviderservices "grovehub/services/mailProvider/services"
	mailprovidervalidators "grovehub/services/mailProvider/validators"
)

type SenderService struct {
	repository    mailproviderrepository.SenderRepository
	domainService mailproviderservices.DomainService
	validator     mailprovidervalidators.SenderValidator
}

func NewSenderService(repository mailproviderrepository.SenderRepository, domainService mailproviderservices.DomainService, validator mailprovidervalidators.SenderValidator) mailproviderservices.SenderService {
	return &SenderService{repository: repository, domainService: domainService, validator: validator}
}

func (service *SenderService) Register(ctx context.Context, input mailproviderservices.RegisterSenderInput) (entity.Sender, error) {
	validated, err := service.validator.ValidateRegister(input)
	if err != nil {
		return entity.Sender{}, err
	}
	domain, err := service.domainService.GetById(ctx, validated.DomainId)
	if err != nil {
		return entity.Sender{}, err
	}
	if domain.GetStatus() != entity.DomainStatusVerified {
		return entity.Sender{}, coreerrors.NewError(mailprovidererror.DomainNotVerified)
	}
	parts := strings.Split(validated.Email, "@")
	if len(parts) != 2 || parts[1] != domain.GetName() {
		return entity.Sender{}, coreerrors.NewError(mailprovidererror.InvalidSender)
	}
	senderId, err := coreid.NewUUIDv7()
	if err != nil {
		return entity.Sender{}, coreerrors.Wrap(mailprovidererror.SenderCreateFailed, err)
	}
	now := time.Now().UTC()
	sender := entity.NewSender(entity.SenderId(senderId), domain.GetDomainId(), validated.Email, validated.DisplayName, entity.SenderStatusActive, 1, now, now)
	if err := service.repository.Create(ctx, sender); err != nil {
		if errors.Is(err, mailproviderrepository.ErrSenderAlreadyExists) {
			return entity.Sender{}, coreerrors.NewError(mailprovidererror.SenderAlreadyExists)
		}
		return entity.Sender{}, coreerrors.Wrap(mailprovidererror.SenderCreateFailed, err)
	}
	return sender, nil
}

func (service *SenderService) GetById(ctx context.Context, senderId string) (entity.Sender, error) {
	validatedId, err := service.validator.ValidateSenderId(senderId)
	if err != nil {
		return entity.Sender{}, err
	}
	sender, err := service.repository.FindById(ctx, validatedId)
	if errors.Is(err, mailproviderrepository.ErrSenderNotFound) {
		return entity.Sender{}, coreerrors.NewError(mailprovidererror.SenderNotFound)
	}
	if err != nil {
		return entity.Sender{}, coreerrors.Wrap(mailprovidererror.SenderReadFailed, err)
	}
	return sender, nil
}

func (service *SenderService) Deactivate(ctx context.Context, senderId string) (entity.Sender, error) {
	sender, err := service.GetById(ctx, senderId)
	if err != nil {
		return entity.Sender{}, err
	}
	if sender.GetStatus() == entity.SenderStatusDeactivated {
		return sender, nil
	}
	updated := sender.SetStatus(entity.SenderStatusDeactivated).SetVersion(sender.GetVersion()+1).SetUpdatedAt(time.Now().UTC())
	if err := service.repository.Update(ctx, updated); err != nil {
		return entity.Sender{}, coreerrors.Wrap(mailprovidererror.SenderUpdateFailed, err)
	}
	return updated, nil
}
