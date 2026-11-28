package impl

import (
	"context"
	"grovehub/libs/schoolgateway"
	"errors"
	"time"
	"grovehub/services/identity/internal/app/entity"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	credential "grovehub/services/identity/internal/app/services/credential"
	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	resetpassword "grovehub/services/identity/internal/app/services/authentication/operations/resetpassword"
)

type Service struct {
	accounts accountrepo.Repository
	sessions sessionrepo.Repository
	tokens tokenrepo.Repository
	credential credential.Service
	unit transaction.UnitOfWork
}

var _ resetpassword.Operation = (*Service)(nil)

func New(accounts accountrepo.Repository, sessions sessionrepo.Repository, tokens tokenrepo.Repository, credentials credential.Service, unit transaction.UnitOfWork) resetpassword.Operation {
	return &Service{accounts: accounts, sessions: sessions, tokens: tokens, credential: credentials, unit: unit}
}

func (op *Service) Execute(ctx context.Context, grant string, password string) error {
	digest, err := secrets.Digest(grant, 32)
	if err != nil {
		return errors.New("invalid reset grant")
	}

	return op.unit.Run(ctx, func(transactionCtx context.Context) error {
		proof, err := op.tokens.Inspect(transactionCtx, digest, entity.ResetGrant, time.Now().UTC())
		if err != nil {
			return err
		}
		account, err := op.accounts.FindByID(transactionCtx, proof.AccountID)
		if err != nil {
			return err
		}
		if account.Status != entity.Active {
			return errors.New("invalid account status")
		}
		if schoolgateway.CheckAccountGeneration(transactionCtx, account.TenantID, proof.SessionGeneration) != nil { return tokenrepo.ErrInvalidToken }
		accountID, err := op.tokens.Consume(transactionCtx, digest, entity.ResetGrant, time.Now().UTC())
		if err != nil { return err }
		if accountID != proof.AccountID { return tokenrepo.ErrInvalidToken }
		if err := op.credential.Replace(transactionCtx, accountID, password); err != nil {
			return err
		}
		_, err = op.sessions.RevokeAll(transactionCtx, accountID)
		return err
	})
}
