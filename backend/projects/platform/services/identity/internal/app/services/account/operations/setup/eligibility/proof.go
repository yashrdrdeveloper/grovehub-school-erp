package eligibility

import (
	"context"
	"grovehub/libs/schoolgateway"
	"time"

	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/accountrole"
	authrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/repository/authority"
	tokens "grovehub/services/identity/internal/app/repository/passwordtoken"
	"grovehub/services/identity/internal/app/repository/rbacmigration"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
)

// ProofGuard resolves scope exclusively from the opaque proof and stored account.
type ProofGuard struct {
	tenant string
	accounts accountrepo.Reader
	tokens tokens.Repository
	locker authority.Locker
	readiness rbacmigration.Readiness
	assignments accountrole.Repository
	snapshots authrepo.Repository
}

func NewProofGuard(tenant string, accounts accountrepo.Reader, tokens tokens.Repository, locker authority.Locker, readiness rbacmigration.Readiness, assignments accountrole.Repository, snapshots authrepo.Repository) *ProofGuard {
	return &ProofGuard{tenant: tenant, accounts: accounts, tokens: tokens, locker: locker, readiness: readiness, assignments: assignments, snapshots: snapshots}
}

func (guard *ProofGuard) Run(ctx context.Context, secret string, purpose entity.TokenPurpose, operation func(context.Context, entity.PasswordToken) error) error {
	if guard == nil || guard.tenant == "" || guard.accounts == nil || guard.tokens == nil || guard.locker == nil || guard.readiness == nil || operation == nil || (purpose != entity.AccountSetup && purpose != entity.SetupGrant) {
		return tokens.ErrInvalidToken
	}
	digest, err := secrets.Digest(secret, 32)
	if err != nil {
		return tokens.ErrInvalidToken
	}
	return guard.locker.Run(ctx, guard.tenant, func(tx context.Context) error {
		proof, err := guard.tokens.Inspect(tx, digest, purpose, time.Now().UTC())
		if err != nil || proof.AccountVersion == nil {
			return tokens.ErrInvalidToken
		}
		account, err := guard.accounts.FindByID(tx, proof.AccountID)
		if err != nil || account.TenantID == nil || *account.TenantID != guard.tenant || account.Version != *proof.AccountVersion {
			return tokens.ErrInvalidToken
		}
		if schoolgateway.CheckTenant(tx, guard.tenant) != nil { return tokens.ErrInvalidToken }
		if purpose == entity.SetupGrant && proof.Attempts >= 2 {
			return tokens.ErrAttemptsExceeded
		}
		if err := guard.readiness.Require(tx, guard.tenant); err != nil {
			return tokens.ErrInvalidToken
		}
		if err := Require(tx, account, guard.assignments, guard.snapshots); err != nil {
			return tokens.ErrInvalidToken
		}
		return operation(tx, proof)
	})
}
