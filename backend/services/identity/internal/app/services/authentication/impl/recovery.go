package impl

import (
	"context"
	"errors"
	"strings"
	"time"

	dto "grovehub/services/identity/internal/app/dto/authentication"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/authority"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	"grovehub/services/identity/internal/app/services/authentication"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	"grovehub/services/identity/internal/app/services/registration"
	"grovehub/services/identity/internal/app/services/security/principal"
)

// RecoveryGuard serializes strict tenant credential operations with invitation
// delivery and normalized account authority changes. Scope is derived from
// stored accounts/proofs, never a tenant supplied by an anonymous request.
type RecoveryGuard struct {
	authentication.Service
	accounts accountrepo.Repository
	tokens tokenrepo.Repository
	sessions sessionrepo.Repository
	locker authority.Locker
	invitations registration.InvitationSender
}

func NewRecoveryGuard(base authentication.Service, accounts accountrepo.Repository, tokens tokenrepo.Repository, sessions sessionrepo.Repository, locker authority.Locker, invitations registration.InvitationSender) authentication.Service {
	return &RecoveryGuard{Service: base, accounts: accounts, tokens: tokens, sessions: sessions, locker: locker, invitations: invitations}
}

func (service *RecoveryGuard) Forgot(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if !strings.Contains(email, "@") { return nil }
	account, err := service.accounts.FindByEmail(ctx, email)
	if err != nil || account.Status != entity.Active || account.TenantID == nil || *account.TenantID == "" { return nil }
	tenant := *account.TenantID
	return service.locker.Run(ctx, tenant, func(tx context.Context) error {
		current, err := service.accounts.FindByID(tx, account.AccountID)
		if err != nil || !eligibleRecoveryAccount(current, tenant) || current.Email != email { return nil }
		return service.invitations.Queue(tx, current, entity.PasswordReset, "there")
	})
}

func eligibleRecoveryAccount(account entity.Account, tenant string) bool {
	return tenant != "" && account.TenantID != nil && *account.TenantID == tenant && account.Status == entity.Active && account.Version > 0
}

func (service *RecoveryGuard) withProof(ctx context.Context, secret string, purpose entity.TokenPurpose, operation func(context.Context, entity.Account) error) error {
	hash, err := secrets.Digest(secret, 32)
	if err != nil { return tokenrepo.ErrInvalidToken }
	proof, err := service.tokens.Inspect(ctx, hash, purpose, time.Now().UTC())
	if err != nil { return tokenrepo.ErrInvalidToken }
	account, err := service.accounts.FindByID(ctx, proof.AccountID)
	if err != nil || account.TenantID == nil || *account.TenantID == "" { return tokenrepo.ErrInvalidToken }
	tenant := *account.TenantID
	return service.locker.Run(ctx, tenant, func(tx context.Context) error {
		currentProof, err := service.tokens.Inspect(tx, hash, purpose, time.Now().UTC())
		if err != nil || currentProof.AccountID != account.AccountID || currentProof.AccountVersion == nil { return tokenrepo.ErrInvalidToken }
		current, err := service.accounts.FindByID(tx, currentProof.AccountID)
		if err != nil || !eligibleRecoveryAccount(current, tenant) || current.Version != *currentProof.AccountVersion { return tokenrepo.ErrInvalidToken }
		return operation(tx, current)
	})
}

func (service *RecoveryGuard) Exchange(ctx context.Context, secret string, purpose string) (authentication.Grant, error) {
	if purpose != "recovery" { return service.Service.Exchange(ctx, secret, purpose) }
	var result authentication.Grant
	err := service.withProof(ctx, secret, entity.PasswordReset, func(tx context.Context, _ entity.Account) error {
		var err error
		result, err = service.Service.Exchange(tx, secret, purpose)
		return err
	})
	if err != nil { return authentication.Grant{}, err }
	return result, nil
}

func (service *RecoveryGuard) InspectProof(ctx context.Context, secret string, purpose string) (dto.GrantResult, error) {
	if purpose != "recovery" && purpose != "reset-grant" { return service.Service.InspectProof(ctx, secret, purpose) }
	tokenPurpose := entity.PasswordReset
	if purpose == "reset-grant" { tokenPurpose = entity.ResetGrant }
	var result dto.GrantResult
	err := service.withProof(ctx, secret, tokenPurpose, func(tx context.Context, _ entity.Account) error {
		var err error
		if tokenPurpose == entity.ResetGrant {
			hash, err := secrets.Digest(secret, 32)
			if err != nil { return tokenrepo.ErrInvalidToken }
			proof, err := service.tokens.Inspect(tx, hash, tokenPurpose, time.Now().UTC())
			if err != nil { return err }
			if proof.Attempts >= 2 { return tokenrepo.ErrResetAttemptsExceeded }
		}
		result, err = service.Service.InspectProof(tx, secret, purpose)
		return err
	})
	if err != nil { return dto.GrantResult{}, err }
	return result, nil
}

func (service *RecoveryGuard) invalidateRecovery(ctx context.Context, accountID string) error {
	if err := service.tokens.Invalidate(ctx, accountID, entity.PasswordReset); err != nil { return err }
	return service.tokens.Invalidate(ctx, accountID, entity.ResetGrant)
}

func (service *RecoveryGuard) ResetPassword(ctx context.Context, secret string, password string) error {
	// Validate tenant and version before consuming an attempt. The actual reset
	// rechecks them under the same authority protocol after reservation.
	if err := service.withProof(ctx, secret, entity.ResetGrant, func(tx context.Context, _ entity.Account) error { return nil }); err != nil {
		return err
	}
	hash, err := secrets.Digest(secret, 32)
	if err != nil { return tokenrepo.ErrInvalidToken }
	if err := service.tokens.ReserveResetAttempt(ctx, hash, time.Now().UTC()); err != nil { return err }
	return service.withProof(ctx, secret, entity.ResetGrant, func(tx context.Context, account entity.Account) error {
		if err := service.Service.ResetPassword(tx, secret, password); err != nil { return err }
		return service.invalidateRecovery(tx, account.AccountID)
	})
}

func (service *RecoveryGuard) Login(ctx context.Context, input dto.Login, userAgent string) (authentication.SessionResult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Identifier))
	account, err := service.accounts.FindByEmail(ctx, email)
	if err != nil || account.Status != entity.Active { return authentication.SessionResult{}, authentication.ErrInvalidCredentials }
	if account.TenantID == nil || *account.TenantID == "" { return service.Service.Login(ctx, input, userAgent) }
	tenant := *account.TenantID
	var result authentication.SessionResult
	err = service.locker.Run(ctx, tenant, func(tx context.Context) error {
		current, err := service.accounts.FindByEmail(tx, email)
		if err != nil || current.AccountID != account.AccountID || !eligibleRecoveryAccount(current, tenant) { return authentication.ErrInvalidCredentials }
		result, err = service.Service.Login(tx, input, userAgent)
		return err
	})
	if err != nil { return authentication.SessionResult{}, err }
	return result, nil
}

func (service *RecoveryGuard) ChangePassword(ctx context.Context, password string) error {
	actor, err := principal.Require(ctx)
	if err != nil { return err }
	if actor.TenantID == "" { return service.Service.ChangePassword(ctx, password) }
	if err := principal.Recent(ctx, 5*time.Minute); err != nil { return err }
	return service.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		account, err := service.accounts.FindByID(tx, actor.AccountID)
		if err != nil || !eligibleRecoveryAccount(account, actor.TenantID) { return principal.ErrAuthenticationRequired }
		now := time.Now().UTC()
		session, err := service.sessions.LockActive(tx, actor.AccountID, actor.SessionID, now)
		if errors.Is(err, sessionrepo.ErrInvalidSession) { return principal.ErrAuthenticationRequired }
		if err != nil { return err }
		if session.ReauthenticatedAt == nil || session.ReauthenticatedAt.After(now) || now.Sub(*session.ReauthenticatedAt) > 5*time.Minute {
			return principal.ErrRecentAuthenticationRequired
		}
		if err := service.Service.ChangePassword(tx, password); err != nil { return err }
		return service.invalidateRecovery(tx, actor.AccountID)
	})
}
