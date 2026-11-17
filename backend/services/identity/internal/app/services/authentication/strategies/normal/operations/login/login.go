package login

import (
	"context"
	"errors"
	"strings"
	"time"

	coreid "grovehub/libs/core/id"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	authdto "grovehub/services/identity/internal/app/dto/authentication"
	"grovehub/services/identity/internal/app/entity"
	auth "grovehub/services/identity/internal/app/services/authentication"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	credential "grovehub/services/identity/internal/app/services/credential"
)

type Operation struct {
	accounts accountrepo.Repository
	sessions sessionrepo.Repository
	credentials credential.Service
}

func New(accounts accountrepo.Repository, sessions sessionrepo.Repository, credentials credential.Service) *Operation {
	return &Operation{accounts: accounts, sessions: sessions, credentials: credentials}
}

func (operation *Operation) Execute(ctx context.Context, input authdto.Login, userAgent string) (auth.SessionResult, error) {
	identifier := strings.ToLower(strings.TrimSpace(input.Identifier))
	if identifier == "" || input.Password == "" || !strings.Contains(identifier, "@") {
		return auth.SessionResult{}, errors.New("invalid credentials")
	}

	account, err := operation.accounts.FindByEmail(ctx, identifier)
	if err != nil || account.Status != entity.Active {
		return auth.SessionResult{}, errors.New("invalid credentials")
	}

	valid, err := operation.credentials.Verify(ctx, account.AccountID, input.Password)
	if err != nil || !valid {
		return auth.SessionResult{}, errors.New("invalid credentials")
	}

	access, err := secrets.NewSecret(32)
	if err != nil {
		return auth.SessionResult{}, err
	}
	refresh, err := secrets.NewSecret(48)
	if err != nil {
		return auth.SessionResult{}, err
	}
	accessHash, err := secrets.Digest(access, 32)
	if err != nil {
		return auth.SessionResult{}, err
	}
	refreshHash, err := secrets.Digest(refresh, 48)
	if err != nil {
		return auth.SessionResult{}, err
	}
	sessionID, err := coreid.NewUUIDv7()
	if err != nil {
		return auth.SessionResult{}, err
	}

	if len(userAgent) > 255 {
		userAgent = userAgent[:255]
	}
	now := time.Now().UTC()
	session := entity.Session{
		SessionID: sessionID,
		AccountID: account.AccountID,
		AccessHash: accessHash,
		RefreshHash: refreshHash,
		AccessExpiresAt: now.Add(15 * time.Minute),
		RefreshExpiresAt: now.Add(30 * 24 * time.Hour),
		ReauthenticatedAt: &now,
		UserAgent: userAgent,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := operation.sessions.Create(ctx, session); err != nil {
		return auth.SessionResult{}, err
	}

	return auth.SessionResult{
		Data: authdto.LoginResult{
			AccountID: account.AccountID,
			PersonGID: account.PersonGID,
			SessionID: sessionID,
			ExpiresIn: 900,
			RequiresSetup: false,
		},
		Cookies: auth.Cookies{Access: access, Refresh: refresh},
	}, nil
}
