package impl

import (
	"context"
	"grovehub/libs/schoolgateway"
	"errors"
	"time"
	"grovehub/services/identity/internal/app/entity"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	auth "grovehub/services/identity/internal/app/services/authentication"
	authdto "grovehub/services/identity/internal/app/dto/authentication"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	exchange "grovehub/services/identity/internal/app/services/authentication/operations/exchange"
)

type Service struct {
	accounts accountrepo.Repository
	tokens tokenrepo.Repository
	unit transaction.UnitOfWork
}

var _ exchange.Operation = (*Service)(nil)

func New(accounts accountrepo.Repository, tokens tokenrepo.Repository, unit transaction.UnitOfWork) exchange.Operation {
	return &Service{accounts: accounts, tokens: tokens, unit: unit}
}

func (op *Service) Execute(ctx context.Context, secret string, purpose string) (auth.Grant, error) {
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

	err = op.unit.Run(ctx, func(transactionCtx context.Context) error {
		proof, err := op.tokens.Inspect(transactionCtx, digest, tokenPurpose, time.Now().UTC())
		if err != nil {
			return err
		}
		account, err := op.accounts.FindByID(transactionCtx, proof.AccountID)
		if err != nil {
			return err
		}
		if tokenPurpose == entity.AccountSetup && (account.Status != entity.PendingSetup || (proof.AccountVersion != nil && *proof.AccountVersion != account.Version)) {
			return errors.New("account setup no longer available")
		}
		if tokenPurpose == entity.PasswordReset && account.Status != entity.Active {
			return errors.New("password reset no longer available")
		}

		if schoolgateway.CheckAccountGeneration(transactionCtx, account.TenantID, proof.SessionGeneration) != nil { return tokenrepo.ErrInvalidToken }
		accountID, err := op.tokens.Consume(transactionCtx, digest, tokenPurpose, time.Now().UTC())
		if err != nil { return err }
		if accountID != proof.AccountID { return tokenrepo.ErrInvalidToken }
		now := time.Now().UTC()
		return op.tokens.Issue(transactionCtx, entity.PasswordToken{
			TokenHash: grantDigest,
			AccountID: accountID,
			Purpose: grantPurpose,
			AccountVersion: proof.AccountVersion,
			SessionGeneration: proof.SessionGeneration,
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
