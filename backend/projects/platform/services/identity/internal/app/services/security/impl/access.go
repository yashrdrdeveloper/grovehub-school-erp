package impl

import (
	"context"
	"errors"
	"time"

	schoolgateway "grovehub/launcher/session"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	security "grovehub/services/identity/internal/app/services/security"
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

func (verifier *Verifier) VerifyAccess(ctx context.Context, accessSecret string) (schoolgateway.Principal, error) {
	hash, err := secrets.Digest(accessSecret, 32)
	if err != nil {
		return schoolgateway.Principal{}, errors.New("invalid access credential")
	}

	session, err := verifier.sessions.FindAccess(ctx, hash, time.Now().UTC())
	if err != nil {
		return schoolgateway.Principal{}, errors.New("session is invalid")
	}

	return verifier.principalForSession(ctx, session)
}

func (verifier *Verifier) VerifyDelegatedSession(ctx context.Context, sessionID string) (schoolgateway.Principal, error) {
 reader, ok := verifier.sessions.(interface {
  FindDelegated(context.Context, string, time.Time) (entity.Session, error)
 })
 if !ok || sessionID == "" { return schoolgateway.Principal{}, errors.New("delegated session unavailable") }
 session, err := reader.FindDelegated(ctx, sessionID, time.Now().UTC())
 if err != nil { return schoolgateway.Principal{}, errors.New("session is invalid") }
 return verifier.principalForSession(ctx, session)
}

func (verifier *Verifier) principalForSession(ctx context.Context, session entity.Session) (schoolgateway.Principal, error) {
	account, err := verifier.accounts.FindByID(ctx, session.AccountID)
	if err != nil || account.Status != entity.Active {
		return schoolgateway.Principal{}, errors.New("account is not active")
	}

	tenantID := ""
	if account.TenantID != nil {
		tenantID = *account.TenantID
	}

	authenticatedAt := session.CreatedAt
	if session.ReauthenticatedAt != nil {
		authenticatedAt = *session.ReauthenticatedAt
	}

	return schoolgateway.Principal{
		AccountID:       account.AccountID,
		PersonGID:       account.PersonGID,
		TenantID:        tenantID,
		SessionID:       session.SessionID,
		SessionGeneration: session.SessionGeneration,
		Admin:           verifier.adminIDs[account.AccountID],
		AuthenticatedAt: authenticatedAt,
	}, nil
}
