// Package impl composes Account use-case operations and exposes the interface.
package impl

import (
	"context"
	"grovehub/libs/schoolgateway"
	"time"
	"grovehub/services/identity/internal/app/entity"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	accountproducer "grovehub/services/identity/internal/app/events/producers/account"
	accountmapper "grovehub/services/identity/internal/app/mappers/account"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	account "grovehub/services/identity/internal/app/services/account"
	search "grovehub/services/identity/internal/app/services/account/operations/search"
	searchimpl "grovehub/services/identity/internal/app/services/account/operations/search/strategies/impl"
	credential "grovehub/services/identity/internal/app/services/credential"
	registration "grovehub/services/identity/internal/app/services/registration"
	accountvalidator "grovehub/services/identity/internal/app/validators/account"
	get "grovehub/services/identity/internal/app/services/account/operations/get"
	getimpl "grovehub/services/identity/internal/app/services/account/operations/get/impl"
	update "grovehub/services/identity/internal/app/services/account/operations/update"
	updateimpl "grovehub/services/identity/internal/app/services/account/operations/update/impl"
	transition "grovehub/services/identity/internal/app/services/account/operations/transition"
	transitionimpl "grovehub/services/identity/internal/app/services/account/operations/transition/impl"
	setup "grovehub/services/identity/internal/app/services/account/operations/setup"
	setupimpl "grovehub/services/identity/internal/app/services/account/operations/setup/impl"
)

// @Service
type Service struct {
	accounts accountrepo.Repository
	tokens tokenrepo.Repository
	registrations registration.Service
	search search.Operation
	get get.Operation
	update update.Operation
	transition transition.Operation
	setup setup.Operation
}

var _ account.Service = (*Service)(nil)

func New(accounts accountrepo.Repository, tokens tokenrepo.Repository, sessions sessionrepo.Repository, credentials credential.Service, registrations registration.Service, mapper accountmapper.Mapper, validator accountvalidator.Validator, unit coredatabase.UnitOfWork, events accountproducer.ChangedProducer, readers ...accountrepo.Reader) account.Service {
	var reader accountrepo.Reader = accounts
	if len(readers) == 1 && readers[0] != nil {
		reader = readers[0]
	}
	return &Service{
		accounts: accounts,
		tokens: tokens,
		registrations: registrations,
		search: search.New(searchimpl.New(reader, mapper, validator)),
		get: getimpl.New(reader, mapper),
		update: updateimpl.New(accounts, mapper, validator, unit, events),
		transition: transitionimpl.New(accounts, sessions, mapper, unit, events),
		setup: setupimpl.New(accounts, tokens, credentials, unit, events),
	}
}

func (service *Service) Create(ctx context.Context, request accountdto.Create, key string) (accountdto.CreateResult, error) {
	return service.registrations.Register(ctx, request, key)
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
	return service.registrations.Resend(ctx, accountID)
}

func (service *Service) Setup(ctx context.Context, grant string, input accountdto.Setup) (accountdto.SetupResult, error) {
	return service.setup.Execute(ctx, grant, input)
}

func (service *Service) SetupStatus(ctx context.Context, grant string) (accountdto.SetupStatusResult, error) {
  digest, err := secrets.Digest(grant, 32)
  if err != nil { return accountdto.SetupStatusResult{}, tokenrepo.ErrInvalidToken }
  now := time.Now().UTC()
  token, err := service.tokens.Inspect(ctx, digest, entity.SetupGrant, now)
  if err != nil { return accountdto.SetupStatusResult{}, err }
  if token.Attempts >= 2 { return accountdto.SetupStatusResult{}, tokenrepo.ErrAttemptsExceeded }
  account, err := service.accounts.FindByID(ctx, token.AccountID)
  if err != nil || account.Status != entity.PendingSetup || schoolgateway.CheckAccountGeneration(ctx, account.TenantID, token.SessionGeneration) != nil { return accountdto.SetupStatusResult{}, tokenrepo.ErrInvalidToken }
  remaining := int(time.Until(token.ExpiresAt).Seconds())
  if remaining < 1 { return accountdto.SetupStatusResult{}, tokenrepo.ErrInvalidToken }
  return accountdto.SetupStatusResult{Verified:true, ExpiresIn:remaining}, nil
}
