package rbac_test

import (
	"context"
	"errors"
	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"os"
	"testing"

	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	roledto "grovehub/services/identity/internal/app/dto/role"
	searchoperation "grovehub/services/identity/internal/app/services/role/operations/search"
	rolevalidator "grovehub/services/identity/internal/app/validators/role"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"grovehub/libs/schoolgateway"
	infrastructure "grovehub/services/infrastructure"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/constants/persona"
	"grovehub/services/identity/internal/app/entity"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	roleimpl "grovehub/services/identity/internal/app/repository/role/impl"
	authorizationrepo "grovehub/services/identity/internal/app/repository/authorization/impl"
	registryrepo "grovehub/services/identity/internal/app/repository/registry"
	registryimpl "grovehub/services/identity/internal/app/repository/registry/impl"
	"grovehub/services/identity/internal/app/services/authorization"
	authorizationimpl "grovehub/services/identity/internal/app/services/authorization/impl"
	"grovehub/services/identity/internal/app/services/registry"
	registryservice "grovehub/services/identity/internal/app/services/registry/impl"
	"grovehub/services/identity/internal/resources/migrations"
)

func TestRBACPostgres(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RBAC PostgreSQL Integration Suite")
}

const actorID = "11111111-1111-4111-8111-111111111111"
const roleID = "22222222-2222-4222-8222-222222222222"
const foreignRoleID = "33333333-3333-4333-8333-333333333333"

var _ = Describe("RBAC persistence against PostgreSQL", func() {
	It("migrates a fresh database, registers inactive capabilities and evaluates live tenant grants", func() {
		dsn := os.Getenv("IDENTITY_RBAC_TEST_DSN")
		if dsn == "" {
			Skip("requires an explicitly provisioned, disposable PostgreSQL database")
		}
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		Expect(err).NotTo(HaveOccurred())
		sqlDB, err := db.DB()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() { Expect(sqlDB.Close()).To(Succeed()) })

		// Never run the historical schema migrations against an existing
		// school database. This integration fixture requires a clean DB.
		var databaseName string
		Expect(db.Raw("SELECT current_database()").Scan(&databaseName).Error).NotTo(HaveOccurred())
		Expect(databaseName).To(Equal("groovehub_rbac_test"))
		var existing int64
		Expect(db.Raw("SELECT count(*) FROM information_schema.schemata WHERE schema_name = 'identity'").Scan(&existing).Error).NotTo(HaveOccurred())
		Expect(existing).To(BeZero(), "integration database must be fresh")
		ctx := context.Background()
		Expect(migrations.Migrate(ctx, sqlDB)).To(Succeed())
		Expect(infrastructure.Migrate(ctx, db)).To(Succeed())

		roles := roleimpl.New(db)
		draft := rolerepo.Draft{TenantID: "school-a", ActorID: actorID, Name: "  Senior Teacher  ", Personas: []persona.Code{persona.Teacher, persona.Coordinator}}
		createdRole, err := roles.CreateDraft(ctx, draft)
		Expect(err).NotTo(HaveOccurred())
		Expect(createdRole.Role.Name).To(Equal("Senior Teacher"))
		Expect(createdRole.Role.Status).To(Equal(entity.RoleInactive))
		Expect(createdRole.Role.Version).To(Equal(int64(1)))
		Expect(createdRole.Personas).To(Equal([]persona.Code{persona.Coordinator, persona.Teacher}))
		_, err = roles.FindByID(ctx, "school-b", createdRole.Role.RoleID)
		Expect(err).To(MatchError(rolerepo.ErrNotFound))
		draft.Name = "senior teacher"
		_, err = roles.CreateDraft(ctx, draft)
		Expect(err).To(MatchError(rolerepo.ErrConflict))
		draft.TenantID = "school-b"
		_, err = roles.CreateDraft(ctx, draft)
		Expect(err).NotTo(HaveOccurred())
		draft.Name = "Invalid role"
		draft.Personas = []persona.Code{persona.Teacher, persona.Teacher}
		_, err = roles.CreateDraft(ctx, draft)
		Expect(err).To(MatchError(rolerepo.ErrInvalidDraft))
		var draftCount int64
		Expect(db.Raw("SELECT count(*) FROM identity.roles").Scan(&draftCount).Error).NotTo(HaveOccurred())
		Expect(draftCount).To(Equal(int64(2)))
		Expect(db.Raw("SELECT count(*) FROM identity.role_permissions").Scan(&draftCount).Error).NotTo(HaveOccurred())
		Expect(draftCount).To(BeZero())

		update := rolerepo.MetadataUpdate{
			Draft: rolerepo.Draft{TenantID: "school-a", ActorID: actorID, Name: "Lead Teacher", Personas: []persona.Code{persona.Teacher}},
			RoleID: createdRole.Role.RoleID, Version: 1,
		}
		updatedRole, err := roles.UpdateMetadata(ctx, update)
		Expect(err).NotTo(HaveOccurred())
		Expect(updatedRole.Role.Version).To(Equal(int64(2)))
		Expect(updatedRole.Role.Status).To(Equal(entity.RoleInactive))
		Expect(updatedRole.Personas).To(Equal([]persona.Code{persona.Teacher}))
		_, err = roles.UpdateMetadata(ctx, update)
		Expect(err).To(MatchError(rolerepo.ErrVersionConflict))
		update.TenantID = "school-b"
		update.Version = 2
		_, err = roles.UpdateMetadata(ctx, update)
		Expect(err).To(MatchError(rolerepo.ErrNotFound))
		update.TenantID = "school-a"
		update.Name = "Concurrent edit"
		type outcome struct { err error }
		outcomes := make(chan outcome, 2)
		for range 2 {
			go func() {
				_, err := roles.UpdateMetadata(ctx, update)
				outcomes <- outcome{err: err}
			}()
		}
		successes := 0
		for range 2 {
			result := <-outcomes
			if result.err == nil {
				successes++
			} else {
				Expect(result.err).To(MatchError(rolerepo.ErrVersionConflict))
			}
		}
		Expect(successes).To(Equal(1))
		updatedRole, err = roles.FindByID(ctx, "school-a", createdRole.Role.RoleID)
		Expect(err).NotTo(HaveOccurred())
		Expect(updatedRole.Role.Version).To(Equal(int64(3)))

		_, err = roles.CreateDraft(ctx, rolerepo.Draft{TenantID: "school-a", ActorID: actorID, Name: "Zulu Teacher", Personas: []persona.Code{persona.Teacher}})
		Expect(err).NotTo(HaveOccurred())
		requestSearch := coresearch.Request{
			PageInfo: coresearch.PageRequest{PageNumber: 1, RowSize: 1},
			Sort: []coresearch.Sort{{Field: roledto.FieldName, Direction: coresearch.DirectionAscending}, {Field: roledto.FieldCreatedAt, Direction: coresearch.DirectionDescending}},
			Filters: []coresearch.Filter{{Field: roledto.FieldPersonas, Operator: coresearch.OperatorEquals, Values: []string{"TEACHER"}}},
		}
		query, err := rolevalidator.NormalizeSearch(requestSearch)
		Expect(err).NotTo(HaveOccurred())
		pageRows, err := roles.Search(ctx, "school-a", query)
		Expect(err).NotTo(HaveOccurred())
		Expect(*pageRows.TotalRecords).To(Equal(int64(2)))
		Expect(pageRows.Items).To(HaveLen(1))
		Expect(pageRows.Items[0].Role.Name).To(Equal("Concurrent edit"))
		Expect(pageRows.Items[0].Personas).To(Equal([]persona.Code{persona.Teacher}))
		page, err := searchoperation.Hydrate(query, pageRows)
		Expect(err).NotTo(HaveOccurred())
		Expect(page.PageInfo.NextCursor).NotTo(BeEmpty())
		requestSearch.PageInfo.PageNumber = 2
		requestSearch.PageInfo.Cursor = page.PageInfo.NextCursor
		query, err = rolevalidator.NormalizeSearch(requestSearch)
		Expect(err).NotTo(HaveOccurred())
		pageRows, err = roles.Search(ctx, "school-a", query)
		Expect(err).NotTo(HaveOccurred())
		Expect(*pageRows.TotalRecords).To(Equal(int64(2)))
		Expect(pageRows.Items).To(HaveLen(1))
		Expect(pageRows.Items[0].Role.Name).To(Equal("Zulu Teacher"))
		Expect(pageRows.HasNext).To(BeFalse())
		requestSearch.PageInfo = coresearch.PageRequest{Rows: 25}
		query, err = rolevalidator.NormalizeSearch(requestSearch)
		Expect(err).NotTo(HaveOccurred())
		pageRows, err = roles.Search(ctx, "school-b", query)
		Expect(err).NotTo(HaveOccurred())
		Expect(*pageRows.TotalRecords).To(Equal(int64(1)))
		Expect(pageRows.Items).To(HaveLen(1))

		catalog := registryservice.New(registryimpl.New(db))
		authorizer := authorizationimpl.New(authorizationrepo.New(db))
		request := registry.Request{
			Features:    []feature.Code{feature.AccountManagement},
			Permissions: []permission.Code{permission.AccountSearch, permission.AccountSelfView},
		}
		result, err := catalog.Register(ctx, request)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(registryrepo.Result{FeaturesCreated: 1, PermissionsCreated: 2}))

		result, err = catalog.Register(ctx, request)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(registryrepo.Result{}))
		var activeCount int64
		Expect(db.Raw("SELECT count(*) FROM identity.permissions WHERE active").Scan(&activeCount).Error).NotTo(HaveOccurred())
		Expect(activeCount).To(BeZero())
		Expect(db.Raw("SELECT count(*) FROM identity.features WHERE active").Scan(&activeCount).Error).NotTo(HaveOccurred())
		Expect(activeCount).To(BeZero())

		// The featureless permission sorts before RoleView; a missing
		// associated feature must roll back the preceding insert as well.
		_, err = catalog.Register(ctx, registry.Request{Permissions: []permission.Code{permission.FeatureView, permission.RoleView}})
		Expect(err).To(MatchError(registryrepo.ErrFeatureMissing))
		var count int64
		Expect(db.Raw("SELECT count(*) FROM identity.permissions WHERE code = ?", permission.FeatureView).Scan(&count).Error).NotTo(HaveOccurred())
		Expect(count).To(BeZero())

		Expect(db.Exec("INSERT INTO identity.accounts(account_id, person_gid, tenant_id, email, status) VALUES (?, ?, ?, ?, 'ACTIVE')", actorID, "44444444-4444-4444-8444-444444444444", "school-a", "actor@example.test").Error).NotTo(HaveOccurred())
		Expect(db.Exec("INSERT INTO identity.roles(role_id, tenant_id, name, normalized_name, status) VALUES (?, 'school-a', 'Manager', 'manager', 'ACTIVE'), (?, 'school-b', 'Manager', 'manager', 'ACTIVE')", roleID, foreignRoleID).Error).NotTo(HaveOccurred())
		Expect(db.Exec("INSERT INTO identity.account_roles(tenant_id, account_id, role_id, assignment_type) VALUES ('school-a', ?, ?, 'PRIMARY')", actorID, roleID).Error).NotTo(HaveOccurred())
		Expect(db.Exec("INSERT INTO identity.account_roles(tenant_id, account_id, role_id, assignment_type) VALUES ('school-a', ?, ?, 'SECONDARY')", actorID, foreignRoleID).Error).To(HaveOccurred())
		Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT 'school-a', ?, permission_id FROM identity.permissions", roleID).Error).NotTo(HaveOccurred())

		mapping, err := roles.Permissions(ctx, "school-a", roleID)
		Expect(err).NotTo(HaveOccurred())
		Expect(mapping.Version).To(Equal(int64(1)))
		Expect(mapping.Codes).To(Equal([]permission.Code{permission.AccountSearch, permission.AccountSelfView}))
		_, err = roles.Permissions(ctx, "school-b", roleID)
		Expect(err).To(MatchError(rolerepo.ErrNotFound))
		emptyMapping, err := roles.Permissions(ctx, "school-a", createdRole.Role.RoleID)
		Expect(err).NotTo(HaveOccurred())
		Expect(emptyMapping.Codes).To(BeEmpty())

		actor := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{AccountID: actorID, TenantID: "school-a", SessionID: "session", Admin: true})
		Expect(authorizer.Require(actor, permission.AccountSearch)).To(MatchError(authorization.ErrForbidden))
		Expect(db.Exec("UPDATE identity.permissions SET active = true").Error).NotTo(HaveOccurred())
		Expect(db.Exec("UPDATE identity.features SET active = true").Error).NotTo(HaveOccurred())
		Expect(authorizer.Require(actor, permission.AccountSelfView)).To(Succeed())
		Expect(authorizer.Require(actor, permission.AccountSearch)).To(MatchError(authorization.ErrForbidden))
		Expect(db.Exec("INSERT INTO identity.tenant_features(tenant_id, feature_id, enabled) SELECT 'school-a', feature_id, true FROM identity.features WHERE code = ?", feature.AccountManagement).Error).NotTo(HaveOccurred())
		Expect(authorizer.Require(actor, permission.AccountSearch)).To(Succeed())

		// Draft mapping writes are deliberately not wired to a service or route.
		writer := roleimpl.NewDraftWriter(db)
		grantUpdate := rolerepo.DraftPermissionUpdate{
			TenantID: "school-a", ActorID: actorID, RoleID: createdRole.Role.RoleID,
			Version: 3, Codes: []permission.Code{permission.AccountSelfView, permission.AccountSearch},
		}
		draftMapping, err := writer.ReplaceDraftPermissions(ctx, grantUpdate)
		Expect(err).NotTo(HaveOccurred())
		Expect(draftMapping.Version).To(Equal(int64(4)))
		Expect(draftMapping.Codes).To(Equal([]permission.Code{permission.AccountSearch, permission.AccountSelfView}))
		_, err = writer.ReplaceDraftPermissions(ctx, grantUpdate)
		Expect(err).To(MatchError(rolerepo.ErrVersionConflict))
		grantUpdate.Version = 4
		sameMapping, err := writer.ReplaceDraftPermissions(ctx, grantUpdate)
		Expect(err).NotTo(HaveOccurred())
		Expect(sameMapping.Version).To(Equal(int64(4)))
		grantUpdate.Codes = []permission.Code{permission.AccountView}
		_, err = writer.ReplaceDraftPermissions(ctx, grantUpdate)
		Expect(err).To(MatchError(rolerepo.ErrInvalidPermissions))
		grantUpdate.Codes = []permission.Code{}
		rollback := errors.New("force rollback")
		Expect(transaction.Run(ctx, db, func(tx context.Context) error {
			cleared, err := writer.ReplaceDraftPermissions(tx, grantUpdate)
			Expect(err).NotTo(HaveOccurred())
			Expect(cleared.Codes).To(BeEmpty())
			Expect(cleared.Version).To(Equal(int64(5)))
			return rollback
		})).To(MatchError(rollback))
		draftMapping, err = roles.Permissions(ctx, "school-a", grantUpdate.RoleID)
		Expect(err).NotTo(HaveOccurred())
		Expect(draftMapping.Version).To(Equal(int64(4)))
		Expect(draftMapping.Codes).To(HaveLen(2))
		grantUpdate.TenantID = "school-b"
		_, err = writer.ReplaceDraftPermissions(ctx, grantUpdate)
		Expect(err).To(MatchError(rolerepo.ErrNotFound))
		grantUpdate.TenantID = "school-a"
		grantUpdate.RoleID = roleID
		grantUpdate.Version = 1
		_, err = writer.ReplaceDraftPermissions(ctx, grantUpdate)
		Expect(err).To(MatchError(rolerepo.ErrLiveRole))
		grantUpdate.RoleID = createdRole.Role.RoleID
		grantUpdate.Version = 4
		Expect(db.Exec("INSERT INTO identity.account_roles(tenant_id, account_id, role_id, assignment_type) VALUES ('school-a', ?, ?, 'SECONDARY')", actorID, grantUpdate.RoleID).Error).NotTo(HaveOccurred())
		_, err = writer.ReplaceDraftPermissions(ctx, grantUpdate)
		Expect(err).To(MatchError(rolerepo.ErrLiveRole))
		Expect(db.Exec("DELETE FROM identity.account_roles WHERE tenant_id = 'school-a' AND role_id = ?", grantUpdate.RoleID).Error).NotTo(HaveOccurred())

		// Re-registration cannot overwrite existing active flags or grants.
		result, err = catalog.Register(ctx, request)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(registryrepo.Result{}))
		Expect(authorizer.Require(actor, permission.AccountSearch)).To(Succeed())
		grants, err := authorizer.Effective(actor)
		Expect(err).NotTo(HaveOccurred())
		Expect(grants).To(Equal([]permission.Code{permission.AccountSearch, permission.AccountSelfView}))

		Expect(db.Exec("UPDATE identity.tenant_features SET enabled = false WHERE tenant_id = 'school-a'").Error).NotTo(HaveOccurred())
		Expect(authorizer.Require(actor, permission.AccountSearch)).To(MatchError(authorization.ErrForbidden))
		Expect(db.Exec("UPDATE identity.roles SET status = 'INACTIVE' WHERE role_id = ?", roleID).Error).NotTo(HaveOccurred())
		Expect(authorizer.Require(actor, permission.AccountSelfView)).To(MatchError(authorization.ErrForbidden))

		// A persisted binding conflict must fail closed instead of being
		// silently repaired or rebound to a different feature.
		Expect(db.Exec("UPDATE identity.permissions SET feature_id = NULL WHERE code = ?", permission.AccountSearch).Error).NotTo(HaveOccurred())
		_, err = catalog.Register(ctx, registry.Request{Permissions: []permission.Code{permission.AccountSearch}})
		Expect(err).To(MatchError(registryrepo.ErrConflict))
		verifyAdministratorGuard(ctx, db)
		verifyApprovalSchema(db)
		verifyApprovalWorkflow(ctx, db)
	})
})
