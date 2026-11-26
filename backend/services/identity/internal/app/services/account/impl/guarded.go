package impl

import (
	"context"

	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	dto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/services/account"
	"grovehub/services/identity/internal/app/services/security/operations/routepolicy"
)

// Guarded protects the same account operations for in-process callers as for
// browser callers. Public setup operations retain their separate grant checks.
type Guarded struct {
	account.Service
	policy *routepolicy.Policy
}

var _ account.Service = (*Guarded)(nil)

func NewGuarded(service account.Service, policy *routepolicy.Policy) account.Service {
	return &Guarded{Service: service, policy: policy}
}

func (service *Guarded) Create(ctx context.Context, request dto.Create, key string) (dto.CreateResult, error) {
	if err := service.policy.Authorize(ctx, "POST /identity/protected/accounts"); err != nil {
		return dto.CreateResult{}, err
	}
	return service.Service.Create(ctx, request, key)
}

func (service *Guarded) Search(ctx context.Context, request coresearch.Request) (coresearch.Result[dto.Account], error) {
	if err := service.policy.Authorize(ctx, "POST /identity/protected/accounts/search"); err != nil {
		return coresearch.Result[dto.Account]{}, err
	}
	return service.Service.Search(ctx, request)
}

func (service *Guarded) Get(ctx context.Context, id string) (dto.Account, error) {
	if err := service.policy.Authorize(ctx, "GET /identity/protected/accounts/{accountId}"); err != nil {
		return dto.Account{}, err
	}
	return service.Service.Get(ctx, id)
}

func (service *Guarded) Update(ctx context.Context, id string, request dto.Update) (dto.Account, error) {
	if err := service.policy.Authorize(ctx, "PATCH /identity/protected/accounts/{accountId}"); err != nil {
		return dto.Account{}, err
	}
	return service.Service.Update(ctx, id, request)
}

func (service *Guarded) Transition(ctx context.Context, id string, action string) (dto.Account, error) {
	if err := service.policy.Authorize(ctx, "PATCH /identity/protected/accounts/{accountId}/"+action); err != nil {
		return dto.Account{}, err
	}
	return service.Service.Transition(ctx, id, action)
}

func (service *Guarded) Resend(ctx context.Context, id string) error {
	if err := service.policy.Authorize(ctx, "POST /identity/protected/accounts/{accountId}/resend-invitation"); err != nil {
		return err
	}
	return service.Service.Resend(ctx, id)
}
