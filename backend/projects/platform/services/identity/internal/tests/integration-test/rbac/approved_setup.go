package rbac_test

import (
	"context"
	"grovehub/libs/schoolgateway"
	"grovehub/services/audit"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"sync/atomic"
	"net/url"
	"time"

	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	"grovehub/services/identity/internal/app/config"
	integrationclient "grovehub/services/integration/client/connection"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	registryrepo "grovehub/services/identity/internal/app/repository/registry/impl"
	"grovehub/services/identity/internal/app/services/registry"
	registryimpl "grovehub/services/identity/internal/app/services/registry/impl"
	"grovehub/services/identity/internal/app/services/security/principal"
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
	policyimpl "grovehub/services/identity/internal/app/repository/approvalpolicy/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	authrepo "grovehub/services/identity/internal/app/repository/authorization/impl"
	credentialrepo "grovehub/services/identity/internal/app/repository/credential/impl"
	mailrepo "grovehub/services/identity/internal/app/repository/invitation/impl"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	tokenimpl "grovehub/services/identity/internal/app/repository/passwordtoken/impl"
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
	sender := invitation.New(tokens, mailrepo.New(db,audit.NewProducer(store)), cfg)
	approvedInvitations := invitation.NewApproved(accounts, assignments, snapshots, sender)
	repository := approvalrepo.New(db, approvalproducer.New(audit.NewProducer(store)))
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
	effects := approvalimpl.Effects{Assignment: replace.New(guard, assignmentimpl.NewWriter(db)), PendingInvitation: approvedInvitations, Audit: resourceproducer.New(audit.NewProducer(store))}
	build := func(effects approvalimpl.Effects) *approvalimpl.Executor {
		return approvalimpl.NewExecutor(repository, locker, guard, readiness, authorizer, snapshots, accounts, roles, assignments, effects)
	}
	fault := effects
	fault.Audit = resourceproducer.New(audit.NewProducer(failingOutbox{}))
	_, err = build(fault).Execute(requesterCtx, approved.RequestID, approved.Version)
	Expect(err).To(MatchError(errInjectedOutbox))
	var count int64
	Expect(db.Raw("SELECT count(*) FROM identity.invitations WHERE account_id = ?", accountID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
	Expect(db.Raw("SELECT count(*) FROM (SELECT event_id, (p#>>'{data,tenantId}') AS tenant_id, (p#>>'{data,accountId}')::uuid AS account_id, (p->>'aggregateId')::uuid AS invitation_id, p#>>'{data,operation}' AS event_type, (p#>>'{data,attempt}')::int AS attempt, p#>>'{data,reason}' AS reason, (p#>>'{data,nextAttemptAt}')::timestamptz AS next_attempt_at, (p->>'occurredAt')::timestamptz AS occurred_at FROM (SELECT event_id,convert_from(payload,'UTF8')::jsonb AS p FROM infrastructure.messaging) m WHERE p->>'eventType' = 'identity.invitation.delivery.changed.v1') a WHERE account_id = ?", accountID).Scan(&count).Error).NotTo(HaveOccurred())
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
	exchangeOperation := exchangeimpl.New(accountimpl.New(db), tokens, coredatabase.NewUnitOfWork(db))
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
	worker := invitation.NewWorker(mailrepo.New(db,audit.NewProducer(store)), accounts, tokens, locker, assignments, snapshots, nil, cfg)
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
	// A legacy generation-zero setup grant cannot cross registry cutover.
	registryCtx := schoolgateway.WithTenantCheck(ctx, func(_ context.Context, id string) (int64,error) {
		Expect(id).To(Equal(tenant))
		return 2,nil
	})
	_, err = onboarding.SetupStatus(registryCtx, grant.Secret)
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))
	_, err = onboarding.Setup(registryCtx, grant.Secret, accountdto.Setup{Password: "Test-setup-password-123"})
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))
	storedGrant, err := tokens.Inspect(ctx, hash, entity.SetupGrant, time.Now().UTC())
	Expect(err).NotTo(HaveOccurred())
	Expect(storedGrant.Attempts).To(BeZero())
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

	// Provision resend capability only in this disposable fixture.
	actor, err := principal.RequireTenant(requesterCtx)
	Expect(err).NotTo(HaveOccurred())
	_, err = registryimpl.New(registryrepo.New(db)).Register(ctx, registry.Request{Permissions: []permission.Code{permission.AccountResendInvitation}})
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.permissions SET active = true WHERE code = ?", permission.AccountResendInvitation).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT ar.tenant_id, ar.role_id, p.permission_id FROM identity.account_roles ar CROSS JOIN identity.permissions p WHERE ar.tenant_id = ? AND ar.account_id = ? AND ar.assignment_type = 'PRIMARY' AND p.code = ? ON CONFLICT DO NOTHING", tenant, actor.AccountID, permission.AccountResendInvitation).Error).NotTo(HaveOccurred())
	Expect(resend.Execute(requesterCtx, accountID)).To(MatchError(accountrepo.ErrConflict))
	Expect(db.Exec("UPDATE identity.invitations SET created_at = ? WHERE id = ?", time.Now().UTC().Add(-time.Minute), mail.ID).Error).NotTo(HaveOccurred())

	// Concurrent authorized resends are serialized: exactly one rotates the proof.
	resends := make(chan error, 2)
	for range 2 {
		go func() { resends <- resend.Execute(requesterCtx, accountID) }()
	}
	successes := 0
	for range 2 {
		if err := <-resends; err == nil {
			successes++
		} else {
			Expect(err).To(MatchError(accountrepo.ErrConflict))
		}
	}
	Expect(successes).To(Equal(1))
	_, err = tokens.Inspect(ctx, hash, entity.SetupGrant, time.Now().UTC())
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))
	var rotated entity.Invitation
	Expect(db.Where("account_id = ? AND id <> ?", accountID, mail.ID).Take(&rotated).Error).NotTo(HaveOccurred())
	rotatedProof := invitationProof(rotated, cfg.EncryptionKey)
	rotatedHash, err := secrets.Digest(rotatedProof, 32)
	Expect(err).NotTo(HaveOccurred())
	rotatedToken, err := tokens.Inspect(ctx, rotatedHash, entity.AccountSetup, time.Now().UTC())
	Expect(err).NotTo(HaveOccurred())
	Expect(rotatedToken.AccountVersion).NotTo(BeNil())
	Expect(*rotatedToken.AccountVersion).To(Equal(before.Version+1))

	// A second authorized resend invalidates a proof that was never exchanged.
	Expect(db.Exec("UPDATE identity.invitations SET created_at = ? WHERE id = ?", time.Now().UTC().Add(-time.Minute), rotated.ID).Error).NotTo(HaveOccurred())
	Expect(resend.Execute(requesterCtx, accountID)).To(Succeed())
	_, err = tokens.Inspect(ctx, rotatedHash, entity.AccountSetup, time.Now().UTC())
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))
	var current entity.Invitation
	Expect(db.Where("account_id = ? AND id NOT IN ?", accountID, []string{mail.ID, rotated.ID}).Take(&current).Error).NotTo(HaveOccurred())
	currentProof := invitationProof(current, cfg.EncryptionKey)

	// Force deterministic claim order without sleeping or touching real transport.
	Expect(db.Exec("UPDATE identity.invitations SET next_attempt_at = ? WHERE id = ?", time.Now().UTC().Add(-time.Minute), rotated.ID).Error).NotTo(HaveOccurred())
	transport := &invitationTransport{}
	transport.fail.Store(true)
	worker = invitation.NewWorker(mailrepo.New(db,audit.NewProducer(store)), accounts, tokens, locker, assignments, snapshots, transport, cfg)
	claimed, err = worker.RunOnce(ctx)
	Expect(err).NotTo(HaveOccurred())
	Expect(claimed).To(BeTrue())
	Expect(transport.calls.Load()).To(BeZero())
	Expect(db.Where("id = ?", rotated.ID).Take(&rotated).Error).NotTo(HaveOccurred())
	Expect(rotated.SuppressedAt).NotTo(BeNil())
	Expect(rotated.SentAt).To(BeNil())

	claimed, err = worker.RunOnce(ctx)
	Expect(claimed).To(BeTrue())
	Expect(err).To(MatchError(errInvitationTransport))
	Expect(transport.calls.Load()).To(Equal(int32(1)))
	Expect(db.Where("id = ?", current.ID).Take(&current).Error).NotTo(HaveOccurred())
	Expect(current.SentAt).To(BeNil())
	Expect(current.LeaseUntil).To(BeNil())
	Expect(current.NextAttemptAt.After(time.Now().UTC())).To(BeTrue())
	Expect(db.Exec("UPDATE identity.invitations SET next_attempt_at = ? WHERE id = ?", time.Now().UTC().Add(-time.Second), current.ID).Error).NotTo(HaveOccurred())
	transport.fail.Store(false)

	// Two workers racing for the single due job may send it only once.
	type deliveryOutcome struct { claimed bool; err error }
	deliveries := make(chan deliveryOutcome, 2)
	for range 2 {
		go func() { claimed, err := worker.RunOnce(ctx); deliveries <- deliveryOutcome{claimed, err} }()
	}
	successes = 0
	for range 2 {
		result := <-deliveries
		Expect(result.err).NotTo(HaveOccurred())
		if result.claimed { successes++ }
	}
	Expect(successes).To(Equal(1))
	Expect(transport.calls.Load()).To(Equal(int32(2)))
	Expect(db.Where("id = ?", current.ID).Take(&current).Error).NotTo(HaveOccurred())
	Expect(current.SentAt).NotTo(BeNil())
	Expect(current.SuppressedAt).To(BeNil())
	Expect(current.Attempts).To(Equal(2))
	verifyInvitationDeliveryAudit(db, tenant, accountID, mail.ID, rotated.ID, current.ID)
	grant, err = enabledAuth.Exchange(ctx, currentProof, "invitation")
	Expect(err).NotTo(HaveOccurred())

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
	verifyRecoveryGuard(ctx, db, tenant, accountID)
}

// The adapter exercises the actual exchange operation without composing unrelated auth use cases.
type setupExchange struct {
 authservice.Service
 operation exchange.Operation
}

func (adapter *setupExchange) Exchange(ctx context.Context, secret string, purpose string) (authservice.Grant, error) {
 return adapter.operation.Execute(ctx, secret, purpose)
}

var errInvitationTransport = errors.New("injected invitation transport failure")

type invitationTransport struct {
	integrationclient.Client
	calls atomic.Int32
	fail atomic.Bool
}

func (transport *invitationTransport) ExecuteSystemMail(ctx context.Context, tenantID string, input json.RawMessage) (any, error) {
	transport.calls.Add(1)
	if transport.fail.Load() {
		return nil, errInvitationTransport
	}
	return nil, nil
}

func invitationProof(mail entity.Invitation, key []byte) string {
	block, err := aes.NewCipher(key)
	Expect(err).NotTo(HaveOccurred())
	aead, err := cipher.NewGCM(block)
	Expect(err).NotTo(HaveOccurred())
	plain, err := aead.Open(nil, mail.Nonce, mail.Ciphertext, []byte(mail.AccountID))
	Expect(err).NotTo(HaveOccurred())
	var content invitation.Encrypted
	Expect(json.Unmarshal(plain, &content)).To(Succeed())
	link, err := url.Parse(content.URL)
	Expect(err).NotTo(HaveOccurred())
	query, err := url.ParseQuery(link.Fragment)
	Expect(err).NotTo(HaveOccurred())
	return query.Get("proof")
}

func verifyInvitationDeliveryAudit(db *gorm.DB, tenant string, accountID string, originalID string, rotatedID string, currentID string) {
	type auditRow struct {
		EventType string
		Attempt int
		Reason *string
		NextAttemptAt *time.Time
	}
	read := func(id string) []auditRow {
		var rows []auditRow
		Expect(db.Raw("SELECT event_type, attempt, reason, next_attempt_at FROM (SELECT event_id, (p#>>'{data,tenantId}') AS tenant_id, (p#>>'{data,accountId}')::uuid AS account_id, (p->>'aggregateId')::uuid AS invitation_id, p#>>'{data,operation}' AS event_type, (p#>>'{data,attempt}')::int AS attempt, p#>>'{data,reason}' AS reason, (p#>>'{data,nextAttemptAt}')::timestamptz AS next_attempt_at, (p->>'occurredAt')::timestamptz AS occurred_at FROM (SELECT event_id,convert_from(payload,'UTF8')::jsonb AS p FROM infrastructure.messaging) m WHERE p->>'eventType' = 'identity.invitation.delivery.changed.v1') a WHERE tenant_id = ? AND account_id = ? AND invitation_id = ? ORDER BY occurred_at,event_id", tenant, accountID, id).Scan(&rows).Error).NotTo(HaveOccurred())
		return rows
	}
	invalid := "INVALID_PROOF"
	for _, id := range []string{originalID, rotatedID} {
		Expect(read(id)).To(Equal([]auditRow{
			{EventType: "QUEUED", Attempt: 0},
			{EventType: "CLAIMED", Attempt: 1},
			{EventType: "SUPPRESSED", Attempt: 1, Reason: &invalid},
		}))
	}
	rows := read(currentID)
	Expect(rows).To(HaveLen(5))
	Expect(rows[0]).To(Equal(auditRow{EventType: "QUEUED", Attempt: 0}))
	Expect(rows[1]).To(Equal(auditRow{EventType: "CLAIMED", Attempt: 1}))
	Expect(rows[2].EventType).To(Equal("RETRY_SCHEDULED"))
	Expect(rows[2].Attempt).To(Equal(1))
	Expect(rows[2].Reason).NotTo(BeNil())
	Expect(*rows[2].Reason).To(Equal("DELIVERY_FAILED"))
	Expect(rows[2].NextAttemptAt).NotTo(BeNil())
	Expect(rows[3]).To(Equal(auditRow{EventType: "CLAIMED", Attempt: 2}))
	Expect(rows[4]).To(Equal(auditRow{EventType: "DELIVERED", Attempt: 2}))
	var count int64
	Expect(db.Raw("SELECT count(*) FROM (SELECT event_id, (p#>>'{data,tenantId}') AS tenant_id, (p#>>'{data,accountId}')::uuid AS account_id, (p->>'aggregateId')::uuid AS invitation_id, p#>>'{data,operation}' AS event_type, (p#>>'{data,attempt}')::int AS attempt, p#>>'{data,reason}' AS reason, (p#>>'{data,nextAttemptAt}')::timestamptz AS next_attempt_at, (p->>'occurredAt')::timestamptz AS occurred_at FROM (SELECT event_id,convert_from(payload,'UTF8')::jsonb AS p FROM infrastructure.messaging) m WHERE p->>'eventType' = 'identity.invitation.delivery.changed.v1') a WHERE account_id = ? AND tenant_id IS DISTINCT FROM ?", accountID, tenant).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
	Expect(read(currentID)).To(Equal(rows))
}
