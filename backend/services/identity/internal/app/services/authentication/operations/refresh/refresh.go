package refresh

import (
	"context"
	"errors"
	"time"

	authdto "grovehub/services/identity/internal/app/dto/authentication"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	auth "grovehub/services/identity/internal/app/services/authentication"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
)

type Operation struct {
	accounts accountrepo.Repository
	sessions sessionrepo.Repository
}

func New(accounts accountrepo.Repository, sessions sessionrepo.Repository) *Operation {
	return &Operation{accounts: accounts, sessions: sessions}
}

func (operation *Operation) Execute(ctx context.Context, token string) (auth.RefreshResult, error) {
	hash, err := secrets.Digest(token, 48)
	if err != nil {
		return auth.RefreshResult{}, errors.New("invalid refresh token")
	}

	session, err := operation.sessions.FindRefresh(ctx, hash, time.Now().UTC())
	if err != nil {
		return auth.RefreshResult{}, errors.New("invalid refresh token")
	}
	account, err := operation.accounts.FindByID(ctx, session.AccountID)
	if err != nil || account.Status != entity.Active {
		return auth.RefreshResult{}, errors.New("invalid refresh token")
	}

	access, err := secrets.NewSecret(32)
	if err != nil {
		return auth.RefreshResult{}, err
	}
	refresh, err := secrets.NewSecret(48)
	if err != nil {
		return auth.RefreshResult{}, err
	}
	accessHash, err := secrets.Digest(access, 32)
	if err != nil {
		return auth.RefreshResult{}, err
	}
	refreshHash, err := secrets.Digest(refresh, 48)
	if err != nil {
		return auth.RefreshResult{}, err
	}

	now := time.Now().UTC()
	rotated := session
	rotated.AccessHash = accessHash
	rotated.RefreshHash = refreshHash
	rotated.AccessExpiresAt = now.Add(15 * time.Minute)
	rotated.UpdatedAt = now
	if err := operation.sessions.Rotate(ctx, rotated, hash); err != nil {
		return auth.RefreshResult{}, errors.New("invalid refresh token")
	}

	return auth.RefreshResult{
		Data: authdto.RefreshResult{
			SessionID: session.SessionID,
			ExpiresIn: 900,
			Refreshed: true,
		},
		Cookies: auth.Cookies{Access: access, Refresh: refresh},
	}, nil
}
