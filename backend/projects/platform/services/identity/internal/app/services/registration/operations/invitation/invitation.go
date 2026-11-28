package invitation

import (
	"context"
	"grovehub/libs/schoolgateway"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"net/url"
	"time"

	"github.com/google/uuid"

	"grovehub/services/identity/internal/app/config"
	"grovehub/services/identity/internal/app/entity"
	registration "grovehub/services/identity/internal/app/services/registration"
	invitationrepo "grovehub/services/identity/internal/app/repository/invitation"
	passwordtokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
)

type Service struct {
	Tokens passwordtokenrepo.Repository
	Mail invitationrepo.Repository
	Config config.Config
}

type Encrypted struct {
	Title string
	Name string
	Description string
	Button string
	URL string
	Minutes int
}

var _ registration.InvitationSender = (*Service)(nil)

func New(tokens passwordtokenrepo.Repository, mail invitationrepo.Repository, config config.Config) registration.InvitationSender {
	return &Service{Tokens: tokens, Mail: mail, Config: config}
}

func (service *Service) Queue(ctx context.Context, account entity.Account, purpose entity.TokenPurpose, name string) error {
	if len(service.Config.EncryptionKey) != 32 {
		return errors.New("invitation encryption unavailable")
	}
	if account.TenantID == nil || *account.TenantID == "" {
		return errors.New("invitation requires a tenant-scoped account")
	}

	generation, err := schoolgateway.TenantGeneration(ctx, *account.TenantID)
	if err != nil { return err }

	target := service.Config.SetupURL
	minutes := 30
	if purpose == entity.PasswordReset {
		target = service.Config.ResetURL
		minutes = 15
	}

	link, err := url.Parse(target)
	if err != nil || link.Scheme != "https" || link.Host == "" {
		return errors.New("invalid trusted reset/setup URL")
	}
	token, err := secrets.NewSecret(32)
	if err != nil {
		return err
	}
	tokenHash, err := secrets.Digest(token, 32)
	if err != nil {
		return err
	}
	if err := service.Tokens.Invalidate(ctx, account.AccountID, purpose); err != nil {
		return err
	}

	var version *int64
	if purpose == entity.PasswordReset {
		if account.Status != entity.Active || account.Version < 1 {
			return errors.New("invalid password reset state")
		}
		version = &account.Version
		if err := service.Tokens.Invalidate(ctx, account.AccountID, entity.ResetGrant); err != nil {
			return err
		}
	}
	if purpose == entity.AccountSetup {
		if account.Status != entity.PendingSetup || account.Version < 1 {
			return errors.New("invalid account setup state")
		}
		version = &account.Version
		if err := service.Tokens.Invalidate(ctx, account.AccountID, entity.SetupGrant); err != nil {
			return err
		}
	}
	now := time.Now().UTC()
	if err := service.Tokens.Issue(ctx, entity.PasswordToken{
		TokenHash: tokenHash,
		AccountID: account.AccountID,
		Purpose: purpose,
		AccountVersion: version,
		SessionGeneration: generation,
		ExpiresAt: now.Add(time.Duration(minutes) * time.Minute),
		CreatedAt: now,
	}); err != nil {
		return err
	}

	// A fragment avoids sending the secret in the initial HTTP GET URL.
	link.Fragment = "proof=" + url.QueryEscape(token)
	mail := Encrypted{
		Title: "Welcome to GrooveHub",
		Name: name,
		Description: "Your GrooveHub school account is ready. Choose your password to finish setup.",
		Button: "Set up your account",
		URL: link.String(),
		Minutes: minutes,
	}
	if purpose == entity.PasswordReset {
		mail.Title = "Reset your GrooveHub password"
		mail.Description = "You requested to reset your GrooveHub password."
		mail.Button = "Reset password"
	}

	plain, err := json.Marshal(mail)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(service.Config.EncryptionKey)
	if err != nil {
		return err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	messageID := uuid.NewString()
	ciphertext := aead.Seal(nil, nonce, plain, []byte(account.AccountID))
	return service.Mail.Queue(ctx, entity.Invitation{
		ID: messageID,
		AccountID: account.AccountID,
		Email: account.Email,
		Purpose: string(purpose),
		Nonce: nonce,
		Ciphertext: ciphertext,
		NextAttemptAt: now,
		CreatedAt: now,
	})
}
