package invitation

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"time"

	integrationclient "grovehub/services/integration/client/connection"
	"grovehub/services/identity/internal/app/config"
	invitationrepo "grovehub/services/identity/internal/app/repository/invitation"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	mailhelpers "grovehub/services/identity/internal/app/services/registration/operations/register/helpers"
)

type Worker struct {
	repository invitationrepo.Repository
	accounts accountrepo.Repository
	client integrationclient.Client
	config config.Config
}

func NewWorker(repository invitationrepo.Repository, accounts accountrepo.Repository, client integrationclient.Client, config config.Config) *Worker {
	return &Worker{repository: repository, accounts: accounts, client: client, config: config}
}

func (worker *Worker) RunOnce(ctx context.Context) (bool, error) {
	invitation, claimed, err := worker.repository.Claim(ctx, time.Now().UTC())
	if err != nil || !claimed {
		return claimed, err
	}
	if len(worker.config.EncryptionKey) != 32 {
		return true, errors.New("mail encryption key unavailable")
	}

	block, err := aes.NewCipher(worker.config.EncryptionKey)
	if err != nil {
		return true, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return true, err
	}
	plain, err := aead.Open(nil, invitation.Nonce, invitation.Ciphertext, []byte(invitation.AccountID))
	if err != nil {
		return true, err
	}

	var data Encrypted
	if err := json.Unmarshal(plain, &data); err != nil {
		return true, err
	}
	mail := mailhelpers.Mail{
		Title: data.Title,
		Name: data.Name,
		Description: data.Description,
		Button: data.Button,
		URL: data.URL,
		Minutes: data.Minutes,
	}

	account, lookupErr := worker.accounts.FindByID(ctx, invitation.AccountID)
	if lookupErr != nil {
		return true, lookupErr
	}
	if account.TenantID == nil || *account.TenantID == "" {
		return true, errors.New("invitation has no tenant scope")
	}
	err = mailhelpers.Send(ctx, worker.client, *account.TenantID, invitation.Email, mail)
	if err != nil {
		delay := time.Duration(1<<min(invitation.Attempts+1, 9)) * time.Second
		_ = worker.repository.Retry(ctx, invitation.ID, time.Now().UTC().Add(delay))
		return true, err
	}
	return true, worker.repository.Delivered(ctx, invitation.ID)
}
