// Package impl composes the registration operations behind the Service interface.
package impl

import (
	"context"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	foundation "grovehub/services/foundation/client/person"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	accountproducer "grovehub/services/identity/internal/app/events/producers/account"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	registration "grovehub/services/identity/internal/app/services/registration"
	register "grovehub/services/identity/internal/app/services/registration/operations/register"
	resend "grovehub/services/identity/internal/app/services/registration/operations/resend"
	resendimpl "grovehub/services/identity/internal/app/services/registration/operations/resend/impl"
	validator "grovehub/services/identity/internal/app/validators/account"
)

// @Service
type Service struct {
	register register.Operation
	resend resend.Operation
}

var _ registration.Service = (*Service)(nil)

func New(accounts accountrepo.Repository, foundationClient foundation.Client, validator validator.Validator, invitations registration.InvitationSender, unit transaction.UnitOfWork, events accountproducer.ChangedProducer) registration.Service {
	return &Service{
		register: register.New(accounts, foundationClient, validator, invitations, unit, events),
		resend: resendimpl.New(accounts, invitations, unit),
	}
}

func (service *Service) Register(ctx context.Context, request accountdto.Create, key string) (accountdto.CreateResult, error) {
	return service.register.Execute(ctx, request, key)
}

func (service *Service) Resend(ctx context.Context, accountID string) error {
	return service.resend.Execute(ctx, accountID)
}
