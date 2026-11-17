// Package impl composes Account use-case operations and exposes the interface.
package impl

import (
	"context"
	transaction "grovehub/libs/core/database/transaction"
	coresearch "grovehub/libs/core/search"
	foundation "grovehub/services/foundation/client/person"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	accountproducer "grovehub/services/identity/internal/app/events/producers/account"
	accountmapper "grovehub/services/identity/internal/app/mappers/account"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	account "grovehub/services/identity/internal/app/services/account"
	create "grovehub/services/identity/internal/app/services/account/operations/create"
	search "grovehub/services/identity/internal/app/services/account/operations/search"
	credential "grovehub/services/identity/internal/app/services/credential"
	registration "grovehub/services/identity/internal/app/services/registration"
	accountvalidator "grovehub/services/identity/internal/app/validators/account"
	get "grovehub/services/identity/internal/app/services/account/operations/get"
	getimpl "grovehub/services/identity/internal/app/services/account/operations/get/impl"
	update "grovehub/services/identity/internal/app/services/account/operations/update"
	updateimpl "grovehub/services/identity/internal/app/services/account/operations/update/impl"
	transition "grovehub/services/identity/internal/app/services/account/operations/transition"
	transitionimpl "grovehub/services/identity/internal/app/services/account/operations/transition/impl"
	resend "grovehub/services/identity/internal/app/services/account/operations/resend"
	resendimpl "grovehub/services/identity/internal/app/services/account/operations/resend/impl"
	setup "grovehub/services/identity/internal/app/services/account/operations/setup"
	setupimpl "grovehub/services/identity/internal/app/services/account/operations/setup/impl"
)

// @Service
type Service struct {
	create *create.Operation
	search *search.Operation
	get get.Operation
	update update.Operation
	transition transition.Operation
	resend resend.Operation
	setup setup.Operation
}

var _ account.Service = (*Service)(nil)

func New(accounts accountrepo.Repository, tokens tokenrepo.Repository, sessions sessionrepo.Repository, foundationClient foundation.Client, credentials credential.Service, invitations registration.Service, mapper accountmapper.Mapper, validator accountvalidator.Validator, unit transaction.UnitOfWork, events accountproducer.ChangedProducer) account.Service {
	return &Service{
		create: create.New(accounts, foundationClient, validator, invitations, unit, events),
		search: search.New(accounts, mapper, validator),
		get: getimpl.New(accounts, mapper),
		update: updateimpl.New(accounts, mapper, validator, unit, events),
		transition: transitionimpl.New(accounts, sessions, mapper, unit, events),
		resend: resendimpl.New(accounts, invitations, unit),
		setup: setupimpl.New(accounts, tokens, credentials, validator, unit, events),
	}
}

func (service *Service) Create(ctx context.Context, request accountdto.Create, key string) (accountdto.CreateResult, error) {
	return service.create.Execute(ctx, request, key)
}

func (service *Service) Search(ctx context.Context, request coresearch.Request) (coresearch.Result[accountdto.Account], error) {
	return service.search.Execute(ctx, request)
}

func (service *Service) Get(ctx context.Context, accountID string) (accountdto.Account, error) {
	return service.get.Execute(ctx, accountID)
}

func (service *Service) Update(ctx context.Context, accountID string, request accountdto.Update) (accountdto.Account, error) {
	return service.update.Execute(ctx, accountID, request)
}

func (service *Service) Transition(ctx context.Context, accountID string, action string) (accountdto.Account, error) {
	return service.transition.Execute(ctx, accountID, action)
}

func (service *Service) Resend(ctx context.Context, accountID string) error {
	return service.resend.Execute(ctx, accountID)
}

func (service *Service) Setup(ctx context.Context, grant string, input accountdto.Setup) (accountdto.SetupResult, error) {
	return service.setup.Execute(ctx, grant, input)
}
