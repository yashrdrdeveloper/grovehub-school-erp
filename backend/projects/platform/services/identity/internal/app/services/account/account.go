package account

import (
	"context"

	coresearch "grovehub/libs/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
)

// Service composes resource operations and presents their stable contract.
type Service interface {
	Create(context.Context, accountdto.Create, string) (accountdto.CreateResult, error)
	Get(context.Context, string) (accountdto.Account, error)
	Search(context.Context, coresearch.Request) (coresearch.Result[accountdto.Account], error)
	Update(context.Context, string, accountdto.Update) (accountdto.Account, error)
	Transition(context.Context, string, string) (accountdto.Account, error)
	Resend(context.Context, string) error
	Setup(context.Context, string, accountdto.Setup) (accountdto.SetupResult, error)
	SetupStatus(context.Context, string) (accountdto.SetupStatusResult, error)
}
