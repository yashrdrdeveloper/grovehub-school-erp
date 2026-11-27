package rbac_test

import (
	"context"
	"grovehub/services/audit"
	"grovehub/services/infrastructure"
	"time"

	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/config"
	"grovehub/services/identity/internal/app/services/security/principal"
	dto "grovehub/services/identity/internal/app/dto/authentication"
	"grovehub/services/identity/internal/app/entity"
	mapper "grovehub/services/identity/internal/app/mappers/account/impl"
	accountrepo "grovehub/services/identity/internal/app/repository/account/impl"
	authorityrepo "grovehub/services/identity/internal/app/repository/authority/impl"
	credentialrepo "grovehub/services/identity/internal/app/repository/credential/impl"
	mailrepo "grovehub/services/identity/internal/app/repository/invitation/impl"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	tokenimpl "grovehub/services/identity/internal/app/repository/passwordtoken/impl"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	sessionimpl "grovehub/services/identity/internal/app/repository/session/impl"
	auth "grovehub/services/identity/internal/app/services/authentication"
	authimpl "grovehub/services/identity/internal/app/services/authentication/impl"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	credentialimpl "grovehub/services/identity/internal/app/services/credential/impl"
	"grovehub/services/identity/internal/app/services/registration/operations/invitation"
)

func verifyRecoveryGuard(ctx context.Context, db *gorm.DB, tenant string, accountID string) {
	accounts := accountrepo.New(db)
	tokens := tokenimpl.New(db)
	sessions := sessionimpl.New(db)
	credentials := credentialimpl.New(credentialrepo.New(db))
	cfg := config.Config{EncryptionKey: make([]byte, 32), ResetURL: "https://school.example/reset-password"}
	sender := invitation.New(tokens, mailrepo.New(db,audit.NewProducer(infrastructure.NewEventStore(db))), cfg)
	base := authimpl.New(accounts, sessions, tokens, credentials, sender, mapper.New(), transaction.NewUnitOfWork(db))
	service := authimpl.NewRecoveryGuard(base, accounts, tokens, sessions, authorityrepo.New(db), sender, audit.NewProducer(infrastructure.NewEventStore(db)))
	account, err := accounts.FindByID(ctx, accountID)
	Expect(err).NotTo(HaveOccurred())
	Expect(account.TenantID).NotTo(BeNil())
	Expect(*account.TenantID).To(Equal(tenant))
	login, err := service.Login(ctx, dto.Login{Identifier: account.Email, Password: "Test-setup-password-123"}, "recovery-test")
	Expect(err).NotTo(HaveOccurred())

	Expect(service.Forgot(ctx, account.Email)).To(Succeed())
	var first entity.Invitation
	Expect(db.Where("account_id = ? AND purpose = ?", accountID, entity.PasswordReset).Take(&first).Error).NotTo(HaveOccurred())
	firstProof := invitationProof(first, cfg.EncryptionKey)
	grant, err := service.Exchange(ctx, firstProof, "recovery")
	Expect(err).NotTo(HaveOccurred())

	// Concurrent forgot requests must leave only one currently valid proof.
	outcomes := make(chan error, 2)
	for range 2 { go func() { outcomes <- service.Forgot(ctx, account.Email) }() }
	for range 2 { Expect(<-outcomes).To(Succeed()) }
	_, err = service.InspectProof(ctx, grant.Secret, "reset-grant")
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))
	var mails []entity.Invitation
	Expect(db.Where("account_id = ? AND purpose = ?", accountID, entity.PasswordReset).Find(&mails).Error).NotTo(HaveOccurred())
	Expect(mails).To(HaveLen(3))
	var current string
	valid := 0
	for _, mail := range mails {
		proof := invitationProof(mail, cfg.EncryptionKey)
		if _, err := service.InspectProof(ctx, proof, "recovery"); err == nil { current = proof; valid++ }
	}
	Expect(valid).To(Equal(1))
	grant, err = service.Exchange(ctx, current, "recovery")
	Expect(err).NotTo(HaveOccurred())

	// Two submissions of the same grant cannot both replace the password.
	for range 2 { go func() { outcomes <- service.ResetPassword(ctx, grant.Secret, "Recovery-password-456") }() }
	successes := 0
	for range 2 {
		if err := <-outcomes; err == nil { successes++ } else { Expect(err).To(MatchError(tokenrepo.ErrInvalidToken)) }
	}
	Expect(successes).To(Equal(1))
	oldAccess, err := secrets.Digest(login.Cookies.Access, 32)
	Expect(err).NotTo(HaveOccurred())
	_, err = sessions.FindAccess(ctx, oldAccess, time.Now().UTC())
	Expect(err).To(MatchError(sessionrepo.ErrInvalidSession))
	_, err = service.Login(ctx, dto.Login{Identifier: account.Email, Password: "Test-setup-password-123"}, "recovery-test")
	Expect(err).To(MatchError(auth.ErrInvalidCredentials))
	currentLogin, err := service.Login(ctx, dto.Login{Identifier: account.Email, Password: "Recovery-password-456"}, "recovery-test")
	Expect(err).NotTo(HaveOccurred())

	// Failed credential replacement rolls back without refunding the attempt.
	Expect(service.Forgot(ctx, account.Email)).To(Succeed())
	var last entity.Invitation
	Expect(db.Where("account_id = ? AND purpose = ?", accountID, entity.PasswordReset).Order("created_at DESC").Take(&last).Error).NotTo(HaveOccurred())
	grant, err = service.Exchange(ctx, invitationProof(last, cfg.EncryptionKey), "recovery")
	Expect(err).NotTo(HaveOccurred())
	grantHash, err := secrets.Digest(grant.Secret, 32)
	Expect(err).NotTo(HaveOccurred())
	for attempt := 1; attempt <= 2; attempt++ {
		Expect(service.ResetPassword(ctx, grant.Secret, "short")).To(HaveOccurred())
		stored, err := tokens.Inspect(ctx, grantHash, entity.ResetGrant, time.Now().UTC())
		Expect(err).NotTo(HaveOccurred())
		Expect(stored.Attempts).To(Equal(attempt))
	}
	_, err = service.InspectProof(ctx, grant.Secret, "reset-grant")
	Expect(err).To(MatchError(tokenrepo.ErrResetAttemptsExceeded))
	Expect(service.ResetPassword(ctx, grant.Secret, "Another-password-789")).To(MatchError(tokenrepo.ErrResetAttemptsExceeded))
	validPassword, err := credentials.Verify(ctx, accountID, "Recovery-password-456")
	Expect(err).NotTo(HaveOccurred())
	Expect(validPassword).To(BeTrue())

	// Exact session validation replaces the bounded session listing.
	authenticated := principal.With(ctx, principal.Principal{TenantID: tenant, AccountID: accountID, SessionID: currentLogin.Data.SessionID, AuthenticatedAt: time.Now().UTC()})
	Expect(service.ChangePassword(authenticated, "Changed-password-789")).To(Succeed())
	Expect(service.ChangePassword(authenticated, "Should-not-change-123")).To(MatchError(principal.ErrAuthenticationRequired))
	_, err = service.InspectProof(ctx, grant.Secret, "reset-grant")
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))
	validPassword, err = credentials.Verify(ctx, accountID, "Changed-password-789")
	Expect(err).NotTo(HaveOccurred())
	Expect(validPassword).To(BeTrue())
	var count int64
	Expect(db.Raw("SELECT count(*) FROM infrastructure.messaging WHERE topic = 'audit_activity_occurred' AND convert_from(payload,'UTF8')::jsonb#>>'{data,assetId}' = ?",accountID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(5)))
	var failed int64
	Expect(db.Raw("SELECT count(*) FROM infrastructure.messaging WHERE topic = 'audit_activity_occurred' AND convert_from(payload,'UTF8')::jsonb#>>'{data,assetId}' = ? AND convert_from(payload,'UTF8')::jsonb#>>'{data,operation}' = 'LOGIN_FAILED' AND convert_from(payload,'UTF8')::jsonb#>>'{data,actorId}' IS NULL",accountID).Scan(&failed).Error).NotTo(HaveOccurred())
	Expect(failed).To(Equal(int64(1)))
}
