package rbac_test

import (
	"context"
	"grovehub/services/audit"

	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	schoolgateway "grovehub/launcher/services/gateway/session"
	foundation "grovehub/services/foundation/client/person"
	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	producer "grovehub/services/identity/internal/app/events/producers/registration/impl"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	accountimpl "grovehub/services/identity/internal/app/repository/account/impl"
	assignmentimpl "grovehub/services/identity/internal/app/repository/accountrole/impl"
	policyimpl "grovehub/services/identity/internal/app/repository/approvalpolicy/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	authrepo "grovehub/services/identity/internal/app/repository/authorization/impl"
	registryrepo "grovehub/services/identity/internal/app/repository/registry/impl"
	"grovehub/services/identity/internal/app/services/authorization"
	authimpl "grovehub/services/identity/internal/app/services/authorization/impl"
	"grovehub/services/identity/internal/app/services/registration/operations/register"
	"grovehub/services/identity/internal/app/services/registry"
	registryimpl "grovehub/services/identity/internal/app/services/registry/impl"
	validator "grovehub/services/identity/internal/app/validators/account/impl"
	messaging "grovehub/services/infrastructure/client/messaging"
)

type registrationFoundation struct {
	foundation.Client
	person foundation.Person
	searches int
}

func (client *registrationFoundation) SearchPersons(ctx context.Context, request foundation.SearchPersonsRequest) (coresearch.Result[foundation.Person], error) {
	client.searches++
	return coresearch.Result[foundation.Person]{Items: []foundation.Person{client.person}}, nil
}

func verifyNormalizedRegistration(ctx context.Context, db *gorm.DB, tenant string, requester string, reviewer string, managementRole string, requiredRole string, store messaging.EventStore) {
	_, err := registryimpl.New(registryrepo.New(db)).Register(ctx, registry.Request{Permissions: []permission.Code{permission.AccountCreate}})
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.permissions SET active = true WHERE code = ?", permission.AccountCreate).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT ?, ?, permission_id FROM identity.permissions WHERE code = ? ON CONFLICT DO NOTHING", tenant, managementRole, permission.AccountCreate).Error).NotTo(HaveOccurred())
	requesterCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: requester, SessionID: "requester-session"})
	reviewerCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: reviewer, SessionID: "reviewer-session"})
	const email = "normalized-registration@approval.test"
	client := &registrationFoundation{person: foundation.Person{GID: uuid.NewString(), Status: "ACTIVE", Contacts: foundation.ResponseContacts{Emails: []foundation.ResponseContactValue{{Value: email, IsPrimary: true}}}}}
	writer := accountimpl.NewRegistrationWriter(db)
	locker := authorityimpl.New(db)
	readiness := policyimpl.NewReadiness(db)
	authorizer := authimpl.New(authrepo.New(db))
	build := func(store messaging.EventStore) register.Operation {
		return register.NewNormalized(writer, client, validator.NewNormalized(), locker, readiness, authorizer, producer.New(audit.NewProducer(store)))
	}
	operation := build(store)
	request := dto.Create{Person: dto.PersonNames{FirstName: "New", LastName: "Account"}, Email: email}
	_, err = operation.Execute(reviewerCtx, request, "registration-1")
	Expect(err).To(MatchError(authorization.ErrForbidden))
	Expect(client.searches).To(BeZero())
	request.RoleAssignment = &dto.RoleAssignment{PrimaryRoleID: managementRole}
	_, err = operation.Execute(requesterCtx, request, "registration-1")
	Expect(err).To(MatchError(accountrepo.ErrConflict))
	Expect(client.searches).To(BeZero())
	request.RoleAssignment = nil
	_, err = build(failingOutbox{}).Execute(requesterCtx, request, "registration-1")
	Expect(err).To(MatchError(errInjectedOutbox))
	var count int64
	Expect(db.Raw("SELECT count(*) FROM identity.accounts WHERE email = ?", email).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
	Expect(db.Raw("SELECT count(*) FROM identity.registration_receipts WHERE tenant_id = ? AND actor_id = ?", tenant, requester).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())

	result, err := operation.Execute(requesterCtx, request, "registration-1")
	Expect(err).NotTo(HaveOccurred())
	Expect(result.Status).To(Equal(string(entity.PendingSetup)))
	Expect(result.Invitation.Status).To(Equal("AWAITING_ROLE_ASSIGNMENT"))
	searches := client.searches
	replayed, err := operation.Execute(requesterCtx, request, "registration-1")
	Expect(err).NotTo(HaveOccurred())
	Expect(replayed).To(Equal(result))
	Expect(client.searches).To(Equal(searches))
	request.Person.FirstName = "Different"
	_, err = operation.Execute(requesterCtx, request, "registration-1")
	Expect(err).To(MatchError(accountrepo.ErrConflict))
	Expect(client.searches).To(Equal(searches))
	assignment, err := assignmentimpl.New(db).Find(ctx, tenant, result.AccountID)
	Expect(err).NotTo(HaveOccurred())
	Expect(assignment.PrimaryRoleID).To(BeEmpty())
	Expect(assignment.SecondaryRoleIDs).To(BeEmpty())
	Expect(assignment.Version).To(Equal(int64(1)))
	Expect(db.Raw("SELECT count(*) FROM identity.accounts WHERE account_id = ? AND primary_role_id IS NULL AND secondary_role_ids = '[]'::jsonb AND registration_idempotency_key IS NULL", result.AccountID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(1)))
	Expect(db.Raw("SELECT count(*) FROM identity.password_tokens WHERE account_id = ?", result.AccountID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
	Expect(db.Raw("SELECT count(*) FROM identity.invitations WHERE account_id = ?", result.AccountID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
	Expect(db.Raw("SELECT count(*) FROM identity.registration_receipts r JOIN infrastructure.messaging m ON m.event_id = r.event_id::text WHERE r.tenant_id = ? AND r.account_id = ? AND convert_from(m.payload,'UTF8')::jsonb->>'eventType' = 'identity.account.registration.changed.v1'", tenant, result.AccountID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(1)))
	Expect(db.Exec("DELETE FROM identity.registration_receipts WHERE tenant_id = ? AND account_id = ?", tenant, result.AccountID).Error).To(HaveOccurred())
	verifyApprovedSetup(ctx, db, tenant, requesterCtx, reviewerCtx, requiredRole, result.AccountID, store)
	verifyAccountMetadata(ctx, db, tenant, requester, reviewer, managementRole, result.AccountID, store)
}
