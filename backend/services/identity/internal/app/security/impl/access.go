package impl

import (
	"context"
	"errors"
	"time"

	coregateway "grovehub/libs/core/gateway"
	security "grovehub/services/identity/internal/app/security"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
)

type Verifier struct {
	accounts accountrepo.Repository
	sessions sessionrepo.Repository
	adminIDs map[string]bool
}

var _ security.AccessVerifier = (*Verifier)(nil)

func New(accounts accountrepo.Repository, sessions sessionrepo.Repository, adminIDs map[string]bool) security.AccessVerifier {
	return &Verifier{accounts: accounts, sessions: sessions, adminIDs: adminIDs}
}

func (verifier *Verifier) VerifyAccess(ctx context.Context, accessSecret string) (coregateway.Principal, error) {
	hash, err := secrets.Digest(accessSecret, 32)
	if err != nil {
		return coregateway.Principal{}, errors.New("invalid access credential")
	}

	session, err := verifier.sessions.FindAccess(ctx, hash, time.Now().UTC())
	if err != nil {
		return coregateway.Principal{}, errors.New("session is invalid")
	}

	account, err := verifier.accounts.FindByID(ctx, session.AccountID)
	if err != nil || account.Status != entity.Active {
		return coregateway.Principal{}, errors.New("account is not active")
	}

	tenantID := ""
	if account.TenantID != nil {
		tenantID = *account.TenantID
	}

	authenticatedAt := session.CreatedAt
	if session.ReauthenticatedAt != nil {
		authenticatedAt = *session.ReauthenticatedAt
	}

	return coregateway.Principal{
		AccountID: account.AccountID,
		PersonGID: account.PersonGID,
		TenantID: tenantID,
		SessionID: session.SessionID,
		Admin: verifier.adminIDs[account.AccountID],
		AuthenticatedAt: authenticatedAt,
	}, nil
}
