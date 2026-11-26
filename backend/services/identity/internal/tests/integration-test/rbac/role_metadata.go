package rbac_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"grovehub/libs/schoolgateway"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/constants/persona"
	controller "grovehub/services/identity/internal/app/controller/role/protected/impl"
	dto "grovehub/services/identity/internal/app/dto/role"
	"grovehub/services/identity/internal/app/entity"
	producer "grovehub/services/identity/internal/app/events/producers/rolemetadata/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	authrepo "grovehub/services/identity/internal/app/repository/authorization/impl"
	capabilityrepo "grovehub/services/identity/internal/app/repository/capability/impl"
	registryrepo "grovehub/services/identity/internal/app/repository/registry/impl"
	repo "grovehub/services/identity/internal/app/repository/role"
	roleimpl "grovehub/services/identity/internal/app/repository/role/impl"
	auditimpl "grovehub/services/identity/internal/app/repository/rolemetadataaudit/impl"
	"grovehub/services/identity/internal/app/services/authorization"
	authimpl "grovehub/services/identity/internal/app/services/authorization/impl"
	capabilityimpl "grovehub/services/identity/internal/app/services/capability/impl"
	"grovehub/services/identity/internal/app/services/registry"
	registryimpl "grovehub/services/identity/internal/app/services/registry/impl"
	roleservice "grovehub/services/identity/internal/app/services/role/impl"
	messaging "grovehub/services/infrastructure/client/messaging"
)

func verifyRoleMetadata(ctx context.Context, db *gorm.DB, tenant string, requester string, reviewer string, managementRole string, store messaging.EventStore) {
	codes := []permission.Code{permission.RoleCreate, permission.RoleUpdate}
	_, err := registryimpl.New(registryrepo.New(db)).Register(ctx, registry.Request{Permissions: codes})
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.permissions SET active = true WHERE code IN ?", codes).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT ?, ?, permission_id FROM identity.permissions WHERE code IN ? ON CONFLICT DO NOTHING", tenant, managementRole, codes).Error).NotTo(HaveOccurred())
	requesterCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: requester, SessionID: "requester-session"})
	reviewerCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: reviewer, SessionID: "reviewer-session"})
	authorizer := authimpl.New(authrepo.New(db))
	locker := authorityimpl.New(db)
	roles := roleimpl.New(db)
	service := roleservice.NewAudited(roles, authorizer, locker, roleimpl.NewDraftInspector(db), producer.New(auditimpl.New(db), store))
	failing := roleservice.NewAudited(roles, authorizer, locker, roleimpl.NewDraftInspector(db), producer.New(auditimpl.New(db), failingOutbox{}))
	draft := dto.CreateDraft{Name: "Audited role draft", Description: "Initial metadata", Personas: []persona.Code{persona.Teacher}}
	_, err = service.CreateDraft(reviewerCtx, draft)
	Expect(err).To(MatchError(authorization.ErrForbidden))
	_, err = failing.CreateDraft(requesterCtx, draft)
	Expect(err).To(MatchError(errInjectedOutbox))
	var count int64
	Expect(db.Raw("SELECT count(*) FROM identity.roles WHERE tenant_id = ? AND normalized_name = 'audited role draft'", tenant).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
	Expect(db.Raw("SELECT count(*) FROM identity.role_metadata_audit WHERE tenant_id = ?", tenant).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())

	created, err := service.CreateDraft(requesterCtx, draft)
	Expect(err).NotTo(HaveOccurred())
	Expect(created.Version).To(Equal(int64(1)))
	Expect(created.Status).To(Equal(string(entity.RoleInactive)))
	draft.Name = " AUDITED ROLE DRAFT "
	_, err = service.CreateDraft(requesterCtx, draft)
	Expect(err).To(MatchError(repo.ErrConflict))
	update := dto.UpdateMetadata{Name: "Audited draft renamed", Description: "Updated metadata", Personas: []persona.Code{persona.Coordinator}, Version: 1}
	_, err = service.UpdateMetadata(requesterCtx, managementRole, update)
	Expect(err).To(MatchError(repo.ErrLiveRole))
	_, err = failing.UpdateMetadata(requesterCtx, created.RoleID, update)
	Expect(err).To(MatchError(errInjectedOutbox))
	preserved, err := roles.FindByID(ctx, tenant, created.RoleID)
	Expect(err).NotTo(HaveOccurred())
	Expect(preserved.Role.Name).To(Equal(created.Name))
	Expect(preserved.Role.Version).To(Equal(int64(1)))
	Expect(preserved.Personas).To(Equal([]persona.Code{persona.Teacher}))

	api := controller.New(service)
	request := httptest.NewRequest(http.MethodPatch, "/identity/protected/roles/"+created.RoleID, strings.NewReader(`{"name":"Audited draft renamed","description":"Updated metadata","personas":["COORDINATOR"],"version":1}`)).WithContext(requesterCtx)
	request.SetPathValue("roleId", created.RoleID)
	response := httptest.NewRecorder()
	api.UpdateMetadata(response, request)
	Expect(response.Code).To(Equal(http.StatusOK))
	updated, err := roles.FindByID(ctx, tenant, created.RoleID)
	Expect(err).NotTo(HaveOccurred())
	Expect(updated.Role.Version).To(Equal(int64(2)))
	Expect(updated.Role.Name).To(Equal(update.Name))
	Expect(updated.Personas).To(Equal(update.Personas))
	_, err = service.UpdateMetadata(requesterCtx, created.RoleID, update)
	Expect(err).To(MatchError(repo.ErrVersionConflict))
	grants, err := roles.Permissions(ctx, tenant, created.RoleID)
	Expect(err).NotTo(HaveOccurred())
	Expect(grants.Codes).To(BeEmpty())

	Expect(db.Raw("SELECT count(*) FROM identity.role_metadata_audit a JOIN infrastructure.messaging m ON m.event_id = a.event_id::text WHERE a.tenant_id = ? AND a.role_id = ? AND a.event_payload = convert_from(m.payload, 'UTF8')::jsonb AND a.event_payload->>'eventType' = 'identity.role.metadata.changed.v1'", tenant, created.RoleID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(2)))
	Expect(db.Raw("SELECT count(*) FROM identity.role_metadata_audit WHERE tenant_id = ? AND role_id = ? AND operation = 'UPDATED' AND event_payload#>>'{data,old,role,name}' = ? AND event_payload#>>'{data,new,role,name}' = ?", tenant, created.RoleID, created.Name, update.Name).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(1)))
	Expect(db.Exec("DELETE FROM identity.role_metadata_audit WHERE tenant_id = ? AND role_id = ?", tenant, created.RoleID).Error).To(HaveOccurred())
	response = httptest.NewRecorder()
	api.CreateDraft(response, httptest.NewRequest(http.MethodPost, "/identity/protected/roles", strings.NewReader(`{"name":"No bypass","personas":["TEACHER"],"status":"ACTIVE"}`)).WithContext(requesterCtx))
	Expect(response.Code).To(Equal(http.StatusBadRequest))
	response = httptest.NewRecorder()
	api.CreateDraft(response, httptest.NewRequest(http.MethodPost, "/identity/protected/roles", strings.NewReader(`{"name":"Another draft","personas":["TEACHER"]}`)).WithContext(requesterCtx))
	Expect(response.Code).To(Equal(http.StatusCreated))
	discovery, err := capabilityimpl.New(capabilityrepo.New(db), locker, authorizer, "").Discover(requesterCtx)
	Expect(err).NotTo(HaveOccurred())
	Expect(discovery.Actions["roleCreate"]).To(BeTrue())
	Expect(discovery.Actions["roleUpdate"]).To(BeTrue())
	Expect(discovery.Actions["accountCreate"]).To(BeFalse())
}
