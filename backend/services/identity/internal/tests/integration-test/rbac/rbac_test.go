package rbac_test

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"grovehub/libs/schoolgateway"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
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

		catalog := registryservice.New(registryimpl.New(db))
		authorizer := authorizationimpl.New(authorizationrepo.New(db))
		request := registry.Request{
			Features: []feature.Code{feature.AccountManagement},
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

		actor := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{AccountID: actorID, TenantID: "school-a", SessionID: "session", Admin: true})
		Expect(authorizer.Require(actor, permission.AccountSearch)).To(MatchError(authorization.ErrForbidden))
		Expect(db.Exec("UPDATE identity.permissions SET active = true").Error).NotTo(HaveOccurred())
		Expect(db.Exec("UPDATE identity.features SET active = true").Error).NotTo(HaveOccurred())
		Expect(authorizer.Require(actor, permission.AccountSelfView)).To(Succeed())
		Expect(authorizer.Require(actor, permission.AccountSearch)).To(MatchError(authorization.ErrForbidden))
		Expect(db.Exec("INSERT INTO identity.tenant_features(tenant_id, feature_id, enabled) SELECT 'school-a', feature_id, true FROM identity.features WHERE code = ?", feature.AccountManagement).Error).NotTo(HaveOccurred())
		Expect(authorizer.Require(actor, permission.AccountSearch)).To(Succeed())

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
	})
})
