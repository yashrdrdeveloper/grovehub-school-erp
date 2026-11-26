package invitation

import (
	"context"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/accountrole"
	authrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/services/account/operations/setup/eligibility"
	"grovehub/services/identity/internal/app/services/registration"
)

type Approved struct {
	accounts accountrepo.Reader
	assignments accountrole.Repository
	snapshots authrepo.Repository
	sender registration.InvitationSender
}

func NewApproved(accounts accountrepo.Reader, assignments accountrole.Repository, snapshots authrepo.Repository, sender registration.InvitationSender) *Approved {
	return &Approved{accounts: accounts, assignments: assignments, snapshots: snapshots, sender: sender}
}

// Queue is invoked only by the approval executor after a successful assignment.
// It joins the same transaction; delivery remains asynchronous.
func (operation *Approved) Queue(ctx context.Context, tenantID string, assignment accountrole.Assignment) error {
	if assignment.Status != entity.PendingSetup {
		return nil
	}
	if !transaction.Active(ctx) || operation.accounts == nil || operation.sender == nil {
		return accountrole.ErrIneligible
	}
	account, err := operation.accounts.FindByID(ctx, assignment.AccountID)
	if err != nil {
		return err
	}
	if account.TenantID == nil || *account.TenantID != tenantID || account.Version != assignment.Version {
		return accountrole.ErrIneligible
	}
	if err := eligibility.Require(ctx, account, operation.assignments, operation.snapshots); err != nil {
		return err
	}
	return operation.sender.Queue(ctx, account, entity.AccountSetup, "there")
}
