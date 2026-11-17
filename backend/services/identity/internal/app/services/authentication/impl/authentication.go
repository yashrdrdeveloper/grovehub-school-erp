package impl

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	transaction "grovehub/libs/core/database/transaction"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	authdto "grovehub/services/identity/internal/app/dto/authentication"
	"grovehub/services/identity/internal/app/entity"
	accountmapper "grovehub/services/identity/internal/app/mappers/account"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	"grovehub/services/identity/internal/app/security/principal"
	auth "grovehub/services/identity/internal/app/services/authentication"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	refresh "grovehub/services/identity/internal/app/services/authentication/operations/refresh"
	strategies "grovehub/services/identity/internal/app/services/authentication/strategies"
	normal "grovehub/services/identity/internal/app/services/authentication/strategies/normal"
	credential "grovehub/services/identity/internal/app/services/credential"
	invitation "grovehub/services/identity/internal/app/services/registration/operations/invitation"
)

// @Service
type Service struct {
	accounts accountrepo.Repository
	sessions sessionrepo.Repository
	tokens tokenrepo.Repository
	credential credential.Service
	invitation *invitation.Service
	mapper accountmapper.Mapper
	unit transaction.UnitOfWork
	login strategies.LoginStrategy
	refresh *refresh.Operation
}

var _ auth.Service = (*Service)(nil)

func New(accounts accountrepo.Repository, sessions sessionrepo.Repository, tokens tokenrepo.Repository, credentials credential.Service, invitations *invitation.Service, mapper accountmapper.Mapper, unit transaction.UnitOfWork) auth.Service {
	return &Service{
		accounts: accounts,
		sessions: sessions,
		tokens: tokens,
		credential: credentials,
		invitation: invitations,
		mapper: mapper,
		unit: unit,
		login: normal.New(accounts, sessions, credentials),
		refresh: refresh.New(accounts, sessions),
	}
}

func (service *Service) Login(ctx context.Context, input authdto.Login, userAgent string) (auth.SessionResult, error) {
	return service.login.Login(ctx, input, userAgent)
}

func (service *Service) Refresh(ctx context.Context, token string) (auth.RefreshResult, error) {
	return service.refresh.Execute(ctx, token)
}

func (service *Service) Forgot(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if !strings.Contains(email, "@") {
		return nil
	}

	account, err := service.accounts.FindByEmail(ctx, email)
	if err != nil || account.Status != entity.Active {
		return nil
	}
	if service.invitation == nil {
		return errors.New("mail unavailable")
	}

	return service.unit.Run(ctx, func(transactionCtx context.Context) error {
		return service.invitation.Queue(transactionCtx, account, entity.PasswordReset, "there")
	})
}

func (service *Service) Exchange(ctx context.Context, secret string, purpose string) (auth.Grant, error) {
	var tokenPurpose entity.TokenPurpose
	var grantPurpose entity.TokenPurpose
	switch purpose {
	case "recovery":
		tokenPurpose = entity.PasswordReset
		grantPurpose = entity.ResetGrant
	case "invitation":
		tokenPurpose = entity.AccountSetup
		grantPurpose = entity.SetupGrant
	default:
		return auth.Grant{}, errors.New("invalid proof purpose")
	}

	digest, err := secrets.Digest(secret, 32)
	if err != nil {
		return auth.Grant{}, errors.New("invalid proof")
	}
	grantSecret, err := secrets.NewSecret(32)
	if err != nil {
		return auth.Grant{}, err
	}
	grantDigest, err := secrets.Digest(grantSecret, 32)
	if err != nil {
		return auth.Grant{}, err
	}

	err = service.unit.Run(ctx, func(transactionCtx context.Context) error {
		accountID, err := service.tokens.Consume(transactionCtx, digest, tokenPurpose, time.Now().UTC())
		if err != nil {
			return err
		}
		account, err := service.accounts.FindByID(transactionCtx, accountID)
		if err != nil {
			return err
		}
		if tokenPurpose == entity.AccountSetup && account.Status != entity.PendingSetup {
			return errors.New("account setup no longer available")
		}
		if tokenPurpose == entity.PasswordReset && account.Status != entity.Active {
			return errors.New("password reset no longer available")
		}

		now := time.Now().UTC()
		return service.tokens.Issue(transactionCtx, entity.PasswordToken{
			TokenHash: grantDigest,
			AccountID: accountID,
			Purpose: grantPurpose,
			ExpiresAt: now.Add(10 * time.Minute),
			CreatedAt: now,
		})
	})
	if err != nil {
		return auth.Grant{}, errors.New("invalid or expired proof")
	}

	return auth.Grant{
		Secret: grantSecret,
		Data: authdto.GrantResult{
			Verified: true,
			Purpose: string(tokenPurpose),
			ExpiresIn: 600,
		},
	}, nil
}

func (service *Service) ResetPassword(ctx context.Context, grant string, password string) error {
	digest, err := secrets.Digest(grant, 32)
	if err != nil {
		return errors.New("invalid reset grant")
	}

	return service.unit.Run(ctx, func(transactionCtx context.Context) error {
		accountID, err := service.tokens.Consume(transactionCtx, digest, entity.ResetGrant, time.Now().UTC())
		if err != nil {
			return err
		}
		account, err := service.accounts.FindByID(transactionCtx, accountID)
		if err != nil {
			return err
		}
		if account.Status != entity.Active {
			return errors.New("invalid account status")
		}
		if err := service.credential.Replace(transactionCtx, accountID, password); err != nil {
			return err
		}
		_, err = service.sessions.RevokeAll(transactionCtx, accountID)
		return err
	})
}

func (service *Service) Reauthenticate(ctx context.Context, password string) error {
	actor, err := principal.Require(ctx)
	if err != nil {
		return err
	}
	valid, err := service.credential.Verify(ctx, actor.AccountID, password)
	if err != nil || !valid {
		return errors.New("invalid credentials")
	}
	return service.sessions.Reauthenticate(ctx, actor.SessionID, time.Now().UTC())
}

func (service *Service) ChangePassword(ctx context.Context, password string) error {
	actor, err := principal.Require(ctx)
	if err != nil {
		return err
	}
	if err := principal.Recent(ctx, 5*time.Minute); err != nil {
		return err
	}

	return service.unit.Run(ctx, func(transactionCtx context.Context) error {
		if err := service.credential.Replace(transactionCtx, actor.AccountID, password); err != nil {
			return err
		}
		_, err := service.sessions.RevokeAll(transactionCtx, actor.AccountID)
		return err
	})
}

func (service *Service) Logout(ctx context.Context) error {
	actor, err := principal.Require(ctx)
	if err != nil {
		return err
	}
	return service.sessions.Revoke(ctx, actor.AccountID, actor.SessionID)
}

func (service *Service) LogoutAll(ctx context.Context) (int64, error) {
	actor, err := principal.Require(ctx)
	if err != nil {
		return 0, err
	}
	if err := principal.Recent(ctx, 5*time.Minute); err != nil {
		return 0, err
	}
	return service.sessions.RevokeAll(ctx, actor.AccountID)
}

func (service *Service) ListSessions(ctx context.Context) (authdto.SessionsResult, error) {
	actor, err := principal.Require(ctx)
	if err != nil {
		return authdto.SessionsResult{}, err
	}
	sessions, err := service.sessions.List(ctx, actor.AccountID)
	if err != nil {
		return authdto.SessionsResult{}, err
	}

	items := make([]authdto.SessionSummary, 0, len(sessions))
	for _, session := range sessions {
		items = append(items, authdto.SessionSummary{
			SessionID: session.SessionID,
			Current: session.SessionID == actor.SessionID,
			CreatedAt: session.CreatedAt,
			LastActiveAt: session.LastActiveAt,
			ExpiresAt: session.RefreshExpiresAt,
		})
	}
	return authdto.SessionsResult{Items: items}, nil
}

func (service *Service) RevokeSession(ctx context.Context, sessionID string) error {
	actor, err := principal.Require(ctx)
	if err != nil {
		return err
	}
	if _, err := uuid.Parse(sessionID); err != nil {
		return errors.New("invalid session ID")
	}
	return service.sessions.Revoke(ctx, actor.AccountID, sessionID)
}

func (service *Service) Me(ctx context.Context) (accountdto.Account, error) {
	actor, err := principal.Require(ctx)
	if err != nil {
		return accountdto.Account{}, err
	}
	account, err := service.accounts.FindByID(ctx, actor.AccountID)
	if err != nil {
		return accountdto.Account{}, err
	}
	return service.mapper.ToResponse(account), nil
}
