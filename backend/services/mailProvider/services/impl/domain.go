package impl

import (
	"context"
	"errors"
	"time"

	coreerrors "grovehub/libs/core/errors"
	coreid "grovehub/libs/core/helpers/id"
	mailprovidererror "grovehub/services/mailProvider/constants/error"
	"grovehub/services/mailProvider/entity"
	mailproviderrepository "grovehub/services/mailProvider/repository"
	mailproviderservices "grovehub/services/mailProvider/services"
	servicehelpers "grovehub/services/mailProvider/services/helpers"
	mailprovidervalidators "grovehub/services/mailProvider/validators"
)

const domainVerificationTimeout = 5 * time.Second

type DomainService struct {
	repository mailproviderrepository.DomainRepository
	validator  mailprovidervalidators.DomainValidator
}

func NewDomainService(
	repository mailproviderrepository.DomainRepository,
	validator mailprovidervalidators.DomainValidator,
) mailproviderservices.DomainService {
	return &DomainService{repository: repository, validator: validator}
}

func (service *DomainService) Register(
	ctx context.Context,
	input mailproviderservices.RegisterDomainInput,
) (entity.Domain, error) {
	validated, err := service.validator.ValidateRegister(input)
	if err != nil {
		return entity.Domain{}, err
	}

	domainId, err := coreid.NewUUIDv7()
	if err != nil {
		return entity.Domain{}, coreerrors.Wrap(mailprovidererror.DomainCreateFailed, err)
	}
	verificationToken, err := coreid.NewUUIDv7()
	if err != nil {
		return entity.Domain{}, coreerrors.Wrap(mailprovidererror.DomainCreateFailed, err)
	}

	now := time.Now().UTC()
	domain := entity.NewDomain(
		entity.DomainId(domainId),
		validated.Name,
		verificationToken,
		entity.DomainStatusPendingVerification,
		1,
		now,
		now,
	)
	if err := service.repository.Create(ctx, domain); err != nil {
		if errors.Is(err, mailproviderrepository.ErrDomainAlreadyExists) {
			return entity.Domain{}, coreerrors.NewError(mailprovidererror.DomainAlreadyExists)
		}
		return entity.Domain{}, coreerrors.Wrap(mailprovidererror.DomainCreateFailed, err)
	}
	return domain, nil
}

func (service *DomainService) GetById(ctx context.Context, domainId string) (entity.Domain, error) {
	validatedId, err := service.validator.ValidateDomainId(domainId)
	if err != nil {
		return entity.Domain{}, err
	}
	domain, err := service.repository.FindById(ctx, validatedId)
	if errors.Is(err, mailproviderrepository.ErrDomainNotFound) {
		return entity.Domain{}, coreerrors.NewError(mailprovidererror.DomainNotFound)
	}
	if err != nil {
		return entity.Domain{}, coreerrors.Wrap(mailprovidererror.DomainReadFailed, err)
	}
	return domain, nil
}

func (service *DomainService) Verify(ctx context.Context, domainId string) (entity.Domain, error) {
	domain, err := service.GetById(ctx, domainId)
	if err != nil {
		return entity.Domain{}, err
	}
	if domain.GetStatus() == entity.DomainStatusVerified {
		return domain, nil
	}
	if domain.GetStatus() == entity.DomainStatusDeactivated {
		return entity.Domain{}, coreerrors.NewError(mailprovidererror.InvalidDomain)
	}

	verificationContext, cancel := context.WithTimeout(ctx, domainVerificationTimeout)
	defer cancel()
	verified, err := servicehelpers.VerifyDomainTXT(
		verificationContext,
		domain.GetName(),
		domain.GetVerificationToken(),
	)
	if err != nil {
		return entity.Domain{}, coreerrors.Wrap(mailprovidererror.DomainVerificationFailed, err)
	}
	if !verified {
		return entity.Domain{}, coreerrors.NewError(mailprovidererror.DomainNotVerified)
	}

	updated := domain.
		SetStatus(entity.DomainStatusVerified).
		SetVersion(domain.GetVersion() + 1).
		SetUpdatedAt(time.Now().UTC())
	if err := service.repository.Update(ctx, updated); err != nil {
		return entity.Domain{}, coreerrors.Wrap(mailprovidererror.DomainUpdateFailed, err)
	}
	return updated, nil
}

func (service *DomainService) Deactivate(ctx context.Context, domainId string) (entity.Domain, error) {
	domain, err := service.GetById(ctx, domainId)
	if err != nil {
		return entity.Domain{}, err
	}
	if domain.GetStatus() == entity.DomainStatusDeactivated {
		return domain, nil
	}

	updated := domain.
		SetStatus(entity.DomainStatusDeactivated).
		SetVersion(domain.GetVersion() + 1).
		SetUpdatedAt(time.Now().UTC())
	if err := service.repository.Update(ctx, updated); err != nil {
		return entity.Domain{}, coreerrors.Wrap(mailprovidererror.DomainUpdateFailed, err)
	}
	return updated, nil
}
