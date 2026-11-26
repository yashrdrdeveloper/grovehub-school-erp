package rbac_test

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"net/url"
	"time"

	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/config"
	accountservice "grovehub/services/identity/internal/app/services/account/impl"
	authservice "grovehub/services/identity/internal/app/services/authentication"
	authruntime "grovehub/services/identity/internal/app/services/authentication/impl"
	"grovehub/services/identity/internal/app/services/approval/operations/cutover"
	resendimpl "grovehub/services/identity/internal/app/services/registration/operations/resend/impl"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/authentication/operations/exchange"
	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	assignmentdto "grovehub/services/identity/internal/app/dto/accountrole"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	accountproducer "grovehub/services/identity/internal/app/events/producers/account/impl"
	approvalproducer "grovehub/services/identity/internal/app/events/producers/approval/impl"
	resourceproducer "grovehub/services/identity/internal/app/events/producers/resource/impl"
	accountimpl "grovehub/services/identity/internal/app/repository/account/impl"
	assignmentimpl "grovehub/services/identity/internal/app/repository/accountrole/impl"
	assignmentrepo "grovehub/services/identity/internal/app/repository/accountrole"
	approvalrepo "grovehub/services/identity/internal/app/repository/approval/impl"
	auditimpl "grovehub/services/identity/internal/app/repository/approvalaudit/impl"
	policyimpl "grovehub/services/identity/internal/app/repository/approvalpolicy/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	authrepo "grovehub/services/identity/internal/app/repository/authorization/impl"
	credentialrepo "grovehub/services/identity/internal/app/repository/credential/impl"
	mailrepo "grovehub/services/identity/internal/app/repository/invitation/impl"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	tokenimpl "grovehub/services/identity/internal/app/repository/passwordtoken/impl"
	resourceaudit "grovehub/services/identity/internal/app/repository/resourceaudit/impl"
	roleimpl "grovehub/services/identity/internal/app/repository/role/impl"
	setupimpl "grovehub/services/identity/internal/app/services/account/operations/setup/impl"
	"grovehub/services/identity/internal/app/services/account/operations/setup/eligibility"
	"grovehub/services/identity/internal/app/services/accountrole/operations/replace"
	approvalimpl "grovehub/services/identity/internal/app/services/approval/impl"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	exchangeimpl "grovehub/services/identity/internal/app/services/authentication/operations/exchange/impl"
	authimpl "grovehub/services/identity/internal/app/services/authorization/impl"
	"grovehub/services/identity/internal/app/services/authorization/operations/protect"
	credentialimpl "grovehub/services/identity/internal/app/services/credential/impl"
	"grovehub/services/identity/internal/app/services/registration/operations/invitation"
	messaging "grovehub/services/infrastructure/client/messaging"
)

func verifyApprovedSetup(ctx context.Context, db *gorm.DB, tenant string, requesterCtx context.Context, reviewerCtx context.Context, roleID string, accountID string, store messaging.EventStore) {
	locker := authorityimpl.New(db)
	guard := protect.New(locker, authorityimpl.NewAdministratorReader(db))
	readiness := policyimpl.NewReadiness(db)
	accounts := accountimpl.NewNormalizedReader(db)
	assignments := assignmentimpl.New(db)
	snapshots := authrepo.New(db)
	authorizer := authimpl.New(snapshots)
	roles := roleimpl.New(db)
	tokens := tokenimpl.New(db)
	credentials := credentialimpl.New(credentialrepo.New(db))
	cfg := config.Config{ApprovalExecutionTenant: tenant, EncryptionKey: make([]byte, 32), SetupURL: "https://school.example/set-password"}
	sender := invitation.New(tokens, mailrepo.New(db), cfg)
	approvedInvitations := invitation.NewApproved(accounts, assignments, snapshots, sender)
	repository := approvalrepo.New(db, approvalproducer.New(auditimpl.New(db), store))
	service := approvalimpl.New(repository, locker, authorizer, snapshots, accounts, roles, assignments)

	before, err := accounts.FindByID(ctx, accountID)
	Expect(err).NotTo(HaveOccurred())
	Expect(eligibility.Require(ctx, before, assignments, snapshots)).To(MatchError(assignmentrepo.ErrIneligible))
	role, err := roles.FindByID(ctx, tenant, roleID)
	Expect(err).NotTo(HaveOccurred())
	pending, err := service.CreateAssignment(requesterCtx, dto.CreateAssignment{Intent: dto.AssignmentIntent{AccountID: accountID, Version: before.Version, PrimaryRoleID: roleID, RoleVersions: []assignmentdto.RoleVersion{{RoleID: roleID, Version: role.Role.Version}}}, Reason: "Initial approved role", IdempotencyKey: "normalized-initial-assignment"})
	Expect(err).NotTo(HaveOccurred())
	approved, err := service.Decide(reviewerCtx, dto.Review{RequestID: pending.RequestID, Version: pending.Version, Decision: kind.Approve})
	Expect(err).NotTo(HaveOccurred())
	effects := approvalimpl.Effects{Assignment: replace.New(guard, assignmentimpl.NewWriter(db)), PendingInvitation: approvedInvitations, Audit: resourceproducer.New(resourceaudit.New(db), store)}
	build := func(effects approvalimpl.Effects) *approvalimpl.Executor {
		return approvalimpl.NewExecutor(repository, locker, guard, readiness, authorizer, snapshots, accounts, roles, assignments, effects)
	}
	fault := effects
	fault.Audit = resourceproducer.New(resourceaudit.New(db), failingOutbox{})
	_, err = build(fault).Execute(requesterCtx, approved.RequestID, approved.Version)
	Expect(err).To(MatchError(errInjectedOutbox))
	var count int64
	Expect(db.Raw("SELECT count(*) FROM identity.invitations WHERE account_id = ?", accountID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
	Expect(db.Raw("SELECT count(*) FROM identity.password_tokens WHERE account_id = ?", accountID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
	executor := build(effects)
	applied, err := executor.Execute(requesterCtx, approved.RequestID, approved.Version)
	Expect(err).NotTo(HaveOccurred())
	Expect(applied.Status).To(Equal(kind.Applied))
	_, err = executor.Execute(requesterCtx, approved.RequestID, approved.Version)
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Raw("SELECT count(*) FROM identity.invitations WHERE account_id = ?", accountID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(1)))

	var mail entity.Invitation
	Expect(db.Where("account_id = ?", accountID).Take(&mail).Error).NotTo(HaveOccurred())
	block, err := aes.NewCipher(cfg.EncryptionKey)
	Expect(err).NotTo(HaveOccurred())
	aead, err := cipher.NewGCM(block)
	Expect(err).NotTo(HaveOccurred())
	plain, err := aead.Open(nil, mail.Nonce, mail.Ciphertext, []byte(accountID))
	Expect(err).NotTo(HaveOccurred())
	var content invitation.Encrypted
	Expect(json.Unmarshal(plain, &content)).To(Succeed())
	link, err := url.Parse(content.URL)
	Expect(err).NotTo(HaveOccurred())
	query, err := url.ParseQuery(link.Fragment)
	Expect(err).NotTo(HaveOccurred())
	proof := query.Get("proof")
	exchangeOperation := exchangeimpl.New(accountimpl.New(db), tokens, transaction.NewUnitOfWork(db))
	proofs := eligibility.NewProofGuard(tenant, accounts, tokens, locker, readiness, assignments, snapshots)
	disabledProofs := eligibility.NewProofGuard("", accounts, tokens, locker, readiness, assignments, snapshots)
	disabledAuth := authruntime.NewOnboarding(&setupExchange{operation: exchangeOperation}, disabledProofs)
	_, err = disabledAuth.Exchange(ctx, proof, "invitation")
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))
	enabledAuth := authruntime.NewOnboarding(&setupExchange{operation: exchangeOperation}, proofs)
	grant, err := enabledAuth.Exchange(ctx, proof, "invitation")
	Expect(err).NotTo(HaveOccurred())
	hash, err := secrets.Digest(grant.Secret, 32)
	Expect(err).NotTo(HaveOccurred())
	token, err := tokens.Inspect(ctx, hash, entity.SetupGrant, time.Now().UTC())
	Expect(err).NotTo(HaveOccurred())
	Expect(token.AccountVersion).NotTo(BeNil())
	Expect(*token.AccountVersion).To(Equal(before.Version+1))

	// The exchanged proof is no longer deliverable, even if the worker never sent it.
	worker := invitation.NewWorker(mailrepo.New(db), accounts, tokens, locker, assignments, snapshots, nil, cfg)
	claimed, err := worker.RunOnce(ctx)
	Expect(err).NotTo(HaveOccurred())
	Expect(claimed).To(BeTrue())
	var suppressed entity.Invitation
	Expect(db.Where("id = ?", mail.ID).Take(&suppressed).Error).NotTo(HaveOccurred())
	Expect(suppressed.SuppressedAt).NotTo(BeNil())
	Expect(suppressed.SuppressionReason).NotTo(BeNil())
	Expect(*suppressed.SuppressionReason).To(Equal("INVALID_PROOF"))
	Expect(suppressed.SentAt).To(BeNil())

	setup := setupimpl.NewNormalized(accounts, accountimpl.NewSetupWriter(db), tokens, credentials, locker, readiness, assignments, snapshots, accountproducer.New(store))
	onboarding := accountservice.NewOnboarding(nil, cutover.New(tenant), nil, nil, setup, proofs)
	status, err := onboarding.SetupStatus(ctx, grant.Secret)
	Expect(err).NotTo(HaveOccurred())
	Expect(status.Verified).To(BeTrue())
	disabled := accountservice.NewOnboarding(nil, cutover.New(""), nil, nil, setup, disabledProofs)
	_, err = disabled.SetupStatus(ctx, grant.Secret)
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))
	_, err = disabled.Setup(ctx, grant.Secret, accountdto.Setup{Password: "Test-setup-password-123"})
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))
	resend := resendimpl.NewNormalized(locker, readiness, authorizer, assignments, approvedInvitations, mailrepo.NewCooldownReader(db))
	Expect(resend.Execute(requesterCtx, accountID)).To(MatchError(authorization.ErrForbidden))
	var recent bool
	Expect(locker.Run(ctx, tenant, func(tx context.Context) error {
		var err error
		recent, err = mailrepo.NewCooldownReader(db).RecentSetup(tx, accountID, time.Now().UTC().Add(-30*time.Second))
		return err
	})).To(Succeed())
	Expect(recent).To(BeTrue())
	stale, err := secrets.NewSecret(32)
	Expect(err).NotTo(HaveOccurred())
	staleHash, err := secrets.Digest(stale, 32)
	Expect(err).NotTo(HaveOccurred())
	now := time.Now().UTC()
	Expect(tokens.Issue(ctx, entity.PasswordToken{TokenHash: staleHash, AccountID: accountID, AccountVersion: &before.Version, Purpose: entity.SetupGrant, CreatedAt: now, ExpiresAt: now.Add(time.Minute)})).To(Succeed())
	_, err = setup.Execute(ctx, stale, accountdto.Setup{Password: "Test-setup-password-123"})
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))

	// Ineligible grants must block activation even after invitation approval.
	Expect(db.Exec("UPDATE identity.permissions SET active = false WHERE code = ?", permission.AccountRoleAssign).Error).NotTo(HaveOccurred())
	_, err = setup.Execute(ctx, grant.Secret, accountdto.Setup{Password: "Test-setup-password-123"})
	Expect(err).To(HaveOccurred())
	Expect(db.Exec("UPDATE identity.permissions SET active = true WHERE code = ?", permission.AccountRoleAssign).Error).NotTo(HaveOccurred())
	Expect(db.Raw("SELECT count(*) FROM identity.password_credentials WHERE account_id = ?", accountID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
	activated, err := onboarding.Setup(ctx, grant.Secret, accountdto.Setup{Password: "Test-setup-password-123"})
	Expect(err).NotTo(HaveOccurred())
	Expect(activated.SetupCompleted).To(BeTrue())
	Expect(activated.Status).To(Equal(string(entity.Active)))
	after, err := accounts.FindByID(ctx, accountID)
	Expect(err).NotTo(HaveOccurred())
	Expect(after.Version).To(Equal(before.Version+2))
	Expect(after.PrimaryRoleID).NotTo(BeNil())
	Expect(*after.PrimaryRoleID).To(Equal(roleID))
	verified, err := credentials.Verify(ctx, accountID, "Test-setup-password-123")
	Expect(err).NotTo(HaveOccurred())
	Expect(verified).To(BeTrue())
	_, err = setup.Execute(ctx, grant.Secret, accountdto.Setup{Password: "Test-setup-password-123"})
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))
}

// The adapter exercises the actual exchange operation without composing unrelated auth use cases.
type setupExchange struct {
 authservice.Service
 operation exchange.Operation
}

func (adapter *setupExchange) Exchange(ctx context.Context, secret string, purpose string) (authservice.Grant, error) {
 return adapter.operation.Execute(ctx, secret, purpose)
}
