// Package impl composes typed Authentication operations; it owns no use-case SQL.
package impl

import (
	"context"
	schoolgateway "grovehub/libs/platform/session"
	"errors"
	"time"
	"grovehub/services/identity/internal/app/entity"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	accountmapper "grovehub/services/identity/internal/app/mappers/account"
	credential "grovehub/services/identity/internal/app/services/credential"
	registration "grovehub/services/identity/internal/app/services/registration"
	auth "grovehub/services/identity/internal/app/services/authentication"
	authdto "grovehub/services/identity/internal/app/dto/authentication"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	strategies "grovehub/services/identity/internal/app/services/authentication/strategies"
	normal "grovehub/services/identity/internal/app/services/authentication/strategies/normal"
	refresh "grovehub/services/identity/internal/app/services/authentication/operations/refresh"
	forgot "grovehub/services/identity/internal/app/services/authentication/operations/forgot"
	forgotimpl "grovehub/services/identity/internal/app/services/authentication/operations/forgot/impl"
	exchange "grovehub/services/identity/internal/app/services/authentication/operations/exchange"
	exchangeimpl "grovehub/services/identity/internal/app/services/authentication/operations/exchange/impl"
	resetpassword "grovehub/services/identity/internal/app/services/authentication/operations/resetpassword"
	resetpasswordimpl "grovehub/services/identity/internal/app/services/authentication/operations/resetpassword/impl"
	reauthenticate "grovehub/services/identity/internal/app/services/authentication/operations/reauthenticate"
	reauthenticateimpl "grovehub/services/identity/internal/app/services/authentication/operations/reauthenticate/impl"
	changepassword "grovehub/services/identity/internal/app/services/authentication/operations/changepassword"
	changepasswordimpl "grovehub/services/identity/internal/app/services/authentication/operations/changepassword/impl"
	logout "grovehub/services/identity/internal/app/services/authentication/operations/logout"
	logoutimpl "grovehub/services/identity/internal/app/services/authentication/operations/logout/impl"
	logoutall "grovehub/services/identity/internal/app/services/authentication/operations/logoutall"
	logoutallimpl "grovehub/services/identity/internal/app/services/authentication/operations/logoutall/impl"
	sessions "grovehub/services/identity/internal/app/services/authentication/operations/sessions"
	sessionsimpl "grovehub/services/identity/internal/app/services/authentication/operations/sessions/impl"
	revokesession "grovehub/services/identity/internal/app/services/authentication/operations/revokesession"
	revokesessionimpl "grovehub/services/identity/internal/app/services/authentication/operations/revokesession/impl"
	me "grovehub/services/identity/internal/app/services/authentication/operations/me"
	meimpl "grovehub/services/identity/internal/app/services/authentication/operations/me/impl"
)

// @Service
type Service struct {
	accounts accountrepo.Repository
	tokens tokenrepo.Repository
	login strategies.LoginStrategy
	refresh refresh.Operation
	forgot forgot.Operation
	exchange exchange.Operation
	resetpassword resetpassword.Operation
	reauthenticate reauthenticate.Operation
	changepassword changepassword.Operation
	logout logout.Operation
	logoutall logoutall.Operation
	sessions sessions.Operation
	revokesession revokesession.Operation
	me me.Operation
}

var _ auth.Service = (*Service)(nil)

func New(accounts accountrepo.Repository, sessions sessionrepo.Repository, tokens tokenrepo.Repository, credentials credential.Service, invitations registration.InvitationSender, mapper accountmapper.Mapper, unit coredatabase.UnitOfWork, readers ...accountrepo.Reader) auth.Service {
	var reader accountrepo.Reader = accounts
	if len(readers) == 1 && readers[0] != nil {
		reader = readers[0]
	}
	return &Service{
		accounts: accounts,
		tokens: tokens,
		login: normal.New(accounts, sessions, credentials),
		refresh: refresh.New(accounts, sessions),
		forgot: forgotimpl.New(accounts, invitations, unit),
		exchange: exchangeimpl.New(accounts, tokens, unit),
		resetpassword: resetpasswordimpl.New(accounts, sessions, tokens, credentials, unit),
		reauthenticate: reauthenticateimpl.New(sessions, credentials),
		changepassword: changepasswordimpl.New(sessions, credentials, unit),
		logout: logoutimpl.New(sessions),
		logoutall: logoutallimpl.New(sessions),
		sessions: sessionsimpl.New(sessions),
		revokesession: revokesessionimpl.New(sessions),
		me: meimpl.New(reader, mapper),
	}
}

func (service *Service) Login(ctx context.Context, input authdto.Login, userAgent string) (auth.SessionResult, error) {
	return service.login.Login(ctx, input, userAgent)
}

func (service *Service) Refresh(ctx context.Context, token string) (auth.RefreshResult, error) {
	return service.refresh.Execute(ctx, token)
}

func (service *Service) Forgot(ctx context.Context, email string) error {
	return service.forgot.Execute(ctx, email)
}

func (service *Service) Exchange(ctx context.Context, secret string, purpose string) (auth.Grant, error) {
	return service.exchange.Execute(ctx, secret, purpose)
}


func (service *Service) InspectProof(ctx context.Context, proof string, purpose string) (authdto.GrantResult, error) {
  var tokenPurpose entity.TokenPurpose
  var expected entity.AccountStatus
  switch purpose {
  case "recovery":
    tokenPurpose, expected = entity.PasswordReset, entity.Active
  case "invitation":
    tokenPurpose, expected = entity.AccountSetup, entity.PendingSetup
  case "reset-grant":
    tokenPurpose, expected = entity.ResetGrant, entity.Active
  default:
    return authdto.GrantResult{}, tokenrepo.ErrInvalidToken
  }
  digest, err := secrets.Digest(proof, 32)
  if err != nil { return authdto.GrantResult{}, tokenrepo.ErrInvalidToken }
  now := time.Now().UTC()
  token, err := service.tokens.Inspect(ctx, digest, tokenPurpose, now)
  if err != nil { return authdto.GrantResult{}, err }
  account, err := service.accounts.FindByID(ctx, token.AccountID)
  if err != nil { return authdto.GrantResult{}, errors.New("invalid or expired proof") }
  if account.Status != expected || schoolgateway.CheckAccountGeneration(ctx, account.TenantID, token.SessionGeneration) != nil { return authdto.GrantResult{}, tokenrepo.ErrInvalidToken }
  remaining := int(time.Until(token.ExpiresAt).Seconds())
  if remaining < 1 { return authdto.GrantResult{}, tokenrepo.ErrInvalidToken }
  return authdto.GrantResult{Verified:true, Purpose:string(tokenPurpose), ExpiresIn:remaining}, nil
}

func (service *Service) ResetPassword(ctx context.Context, grant string, password string) error {
	return service.resetpassword.Execute(ctx, grant, password)
}

func (service *Service) Reauthenticate(ctx context.Context, password string) error {
	return service.reauthenticate.Execute(ctx, password)
}

func (service *Service) ChangePassword(ctx context.Context, password string) error {
	return service.changepassword.Execute(ctx, password)
}

func (service *Service) Logout(ctx context.Context) error {
	return service.logout.Execute(ctx)
}

func (service *Service) LogoutAll(ctx context.Context) (int64, error) {
	return service.logoutall.Execute(ctx)
}

func (service *Service) ListSessions(ctx context.Context) (authdto.SessionsResult, error) {
	return service.sessions.Execute(ctx)
}

func (service *Service) RevokeSession(ctx context.Context, sessionID string) error {
	return service.revokesession.Execute(ctx, sessionID)
}

func (service *Service) Me(ctx context.Context) (accountdto.Account, error) {
	return service.me.Execute(ctx)
}
