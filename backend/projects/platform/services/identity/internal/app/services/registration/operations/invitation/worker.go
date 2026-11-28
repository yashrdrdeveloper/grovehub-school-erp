package invitation

import (
	"context"
	schoolgateway "grovehub/launcher/session"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"net/url"
	"time"

	integrationclient "grovehub/services/integration/client/connection"
	"grovehub/services/identity/internal/app/config"
	"grovehub/services/identity/internal/app/entity"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/accountrole"
	authrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/repository/authority"
	invitationrepo "grovehub/services/identity/internal/app/repository/invitation"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	"grovehub/services/identity/internal/app/services/account/operations/setup/eligibility"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	mailhelpers "grovehub/services/identity/internal/app/services/registration/operations/register/helpers"
)

type Worker struct {
	repository invitationrepo.Repository
	accounts accountrepo.Reader
	tokens tokenrepo.Repository
	locker authority.Locker
	assignments accountrole.Repository
	snapshots authrepo.Repository
	client integrationclient.Client
	config config.Config
	tenantState schoolgateway.TenantStateLookup
}

func NewWorker(repository invitationrepo.Repository, accounts accountrepo.Reader, tokens tokenrepo.Repository, locker authority.Locker, assignments accountrole.Repository, snapshots authrepo.Repository, client integrationclient.Client, config config.Config, lookups ...schoolgateway.TenantStateLookup) *Worker {
	worker := &Worker{repository: repository, accounts: accounts, tokens: tokens, locker: locker, assignments: assignments, snapshots: snapshots, client: client, config: config}
	if len(lookups) == 1 { worker.tenantState = lookups[0] }
	return worker
}

func (worker *Worker) RunOnce(ctx context.Context) (bool, error) {
	invitation, claimed, err := worker.repository.Claim(ctx, time.Now().UTC())
	if err != nil || !claimed {
		return claimed, err
	}
	data, err := worker.decrypt(invitation)
	if err != nil {
		return true, err
	}
	account, err := worker.accounts.FindByID(ctx, invitation.AccountID)
	if errors.Is(err, accountrepo.ErrNotFound) {
		return true, worker.repository.Suppress(ctx, invitation, "ACCOUNT_INELIGIBLE")
	}
	if err != nil {
		return true, err
	}
	if account.TenantID == nil || *account.TenantID == "" {
		return true, worker.repository.Suppress(ctx, invitation, "ACCOUNT_INELIGIBLE")
	}
	tenant := *account.TenantID
	var sendErr error
	err = worker.locker.Run(ctx, tenant, func(tx context.Context) error {
		owned, err := worker.repository.LockClaim(tx, invitation, time.Now().UTC())
		if err != nil || !owned {
			return err
		}
		reason, err := worker.suppressionReason(tx, tenant, invitation, data)
		if err != nil {
			return err
		}
		if reason != "" {
			return worker.repository.Suppress(tx, invitation, reason)
		}

		mail := mailhelpers.Mail{DeliveryID: invitation.ID, Title: data.Title, Name: data.Name, Description: data.Description, Button: data.Button, URL: data.URL, Minutes: data.Minutes}
		sendCtx, cancel := context.WithTimeout(tx, 10*time.Second)
		defer cancel()
		sendErr = mailhelpers.Send(sendCtx, worker.client, tenant, invitation.Email, mail)
		if errors.Is(sendErr, integrationclient.ErrMailDeliveryUncertain) || errors.Is(sendErr, integrationclient.ErrMailDeliveryConflict) {
			return worker.repository.Suppress(tx, invitation, "DELIVERY_UNCERTAIN")
		}
		if sendErr != nil {
			delay := time.Duration(1<<min(invitation.Attempts+1, 9)) * time.Second
			return worker.repository.Retry(tx, invitation.ID, time.Now().UTC().Add(delay))
		}
		return worker.repository.Delivered(tx, invitation.ID)
	})
	if err != nil {
		return true, err
	}
	return true, sendErr
}

func (worker *Worker) decrypt(invitation entity.Invitation) (Encrypted, error) {
	if len(worker.config.EncryptionKey) != 32 {
		return Encrypted{}, errors.New("mail encryption key unavailable")
	}
	block, err := aes.NewCipher(worker.config.EncryptionKey)
	if err != nil {
		return Encrypted{}, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return Encrypted{}, err
	}
	if len(invitation.Nonce) != aead.NonceSize() {
		return Encrypted{}, errors.New("invalid invitation nonce")
	}
	plain, err := aead.Open(nil, invitation.Nonce, invitation.Ciphertext, []byte(invitation.AccountID))
	if err != nil {
		return Encrypted{}, err
	}
	var data Encrypted
	if err := json.Unmarshal(plain, &data); err != nil {
		return Encrypted{}, errors.New("invalid invitation content")
	}
	return data, nil
}

func (worker *Worker) suppressionReason(ctx context.Context, tenant string, invitation entity.Invitation, data Encrypted) (string, error) {
	account, err := worker.accounts.FindByID(ctx, invitation.AccountID)
	if errors.Is(err, accountrepo.ErrNotFound) {
		return "ACCOUNT_INELIGIBLE", nil
	}
	if err != nil {
		return "", err
	}
	if account.TenantID == nil || *account.TenantID != tenant || account.Email != invitation.Email {
		return "ACCOUNT_INELIGIBLE", nil
	}
	purpose := entity.TokenPurpose(invitation.Purpose)
	if purpose != entity.AccountSetup && purpose != entity.PasswordReset {
		return "INVALID_PROOF", nil
	}
	link, err := url.Parse(data.URL)
	if err != nil {
		return "INVALID_PROOF", nil
	}
	query, err := url.ParseQuery(link.Fragment)
	if err != nil {
		return "INVALID_PROOF", nil
	}
	hash, err := secrets.Digest(query.Get("proof"), 32)
	if err != nil {
		return "INVALID_PROOF", nil
	}
	proof, err := worker.tokens.Inspect(ctx, hash, purpose, time.Now().UTC())
	if errors.Is(err, tokenrepo.ErrInvalidToken) {
		return "INVALID_PROOF", nil
	}
	if err != nil {
		return "", err
	}
	if proof.AccountID != account.AccountID {
		return "INVALID_PROOF", nil
	}
	if worker.tenantState != nil {
		state, err := worker.tenantState(ctx, tenant)
		reason, err := tenantDeliveryDecision(state, err, proof.SessionGeneration)
		if err != nil || reason != "" { return reason, err }
	}
	if purpose == entity.PasswordReset {
		if account.Status != entity.Active || (!worker.config.LegacyRBACDevelopment && (proof.AccountVersion == nil || *proof.AccountVersion != account.Version)) {
			return "ACCOUNT_INELIGIBLE", nil
		}
		return "", nil
	}
	if account.Status != entity.PendingSetup || proof.AccountVersion == nil || *proof.AccountVersion != account.Version {
		return "ACCOUNT_INELIGIBLE", nil
	}
	if !worker.config.LegacyRBACDevelopment {
		if worker.config.ApprovalExecutionTenant != tenant {
			return "ACCOUNT_INELIGIBLE", nil
		}
		if err := eligibility.Require(ctx, account, worker.assignments, worker.snapshots); err != nil {
			if errors.Is(err, accountrole.ErrIneligible) {
				return "ACCOUNT_INELIGIBLE", nil
			}
			return "", err
		}
	}
	return "", nil
}
