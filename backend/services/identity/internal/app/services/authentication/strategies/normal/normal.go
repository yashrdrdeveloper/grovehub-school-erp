package normal

import (
	"context"

	"grovehub/services/identity/internal/app/dto/authentication"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	auth "grovehub/services/identity/internal/app/services/authentication"
	login "grovehub/services/identity/internal/app/services/authentication/strategies/normal/operations/login"
	"grovehub/services/identity/internal/app/services/authentication/strategies"
	credential "grovehub/services/identity/internal/app/services/credential"
)

// @Strategy
type Strategy struct {
	login *login.Operation
}

var _ strategies.LoginStrategy = (*Strategy)(nil)

func New(accounts accountrepo.Repository, sessions sessionrepo.Repository, credentials credential.Service) strategies.LoginStrategy {
	return &Strategy{login: login.New(accounts, sessions, credentials)}
}

func (strategy *Strategy) Login(ctx context.Context, input authentication.Login, userAgent string) (auth.SessionResult, error) {
	return strategy.login.Execute(ctx, input, userAgent)
}
