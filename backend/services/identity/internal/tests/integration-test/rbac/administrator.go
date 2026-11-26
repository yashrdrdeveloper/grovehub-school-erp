package rbac_test

import (
	"context"

	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/repository/authority"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	registryimpl "grovehub/services/identity/internal/app/repository/registry/impl"
	"grovehub/services/identity/internal/app/services/authorization/operations/protect"
	"grovehub/services/identity/internal/app/services/registry"
	registryservice "grovehub/services/identity/internal/app/services/registry/impl"
)

func verifyAdministratorGuard(ctx context.Context, db *gorm.DB) {
	const tenant = "guard-school"
	const first = "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	const second = "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
	const grants = "cccccccc-cccc-4ccc-8ccc-cccccccccccc"
	const assignments = "dddddddd-dddd-4ddd-8ddd-dddddddddddd"
	catalog := registryservice.New(registryimpl.New(db))
	_, err := catalog.Register(ctx, registry.Request{Features: []feature.Code{feature.RoleManagement}, Permissions: permission.AdministratorRecoveryCapabilities()})
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.features SET active = true WHERE code = ?", feature.RoleManagement).Error).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.permissions SET active = true WHERE code IN ?", permission.AdministratorRecoveryCapabilities()).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.tenant_features(tenant_id, feature_id, enabled) SELECT ?, feature_id, true FROM identity.features WHERE code = ?", tenant, feature.RoleManagement).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.accounts(account_id, person_gid, tenant_id, email, status) VALUES (?, ?, ?, 'first@guard.test', 'ACTIVE'), (?, ?, ?, 'second@guard.test', 'ACTIVE')", first, first, tenant, second, second, tenant).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.roles(role_id, tenant_id, name, normalized_name, status) VALUES (?, ?, 'Grants', 'grants', 'ACTIVE'), (?, ?, 'Assignments', 'assignments', 'ACTIVE')", grants, tenant, assignments, tenant).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT ?, ?, permission_id FROM identity.permissions WHERE code = ?", tenant, grants, permission.RolePermissionManage).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT ?, ?, permission_id FROM identity.permissions WHERE code = ?", tenant, assignments, permission.AccountRoleAssign).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.account_roles(tenant_id, account_id, role_id, assignment_type) VALUES (?, ?, ?, 'PRIMARY'), (?, ?, ?, 'PRIMARY')", tenant, first, grants, tenant, second, assignments).Error).NotTo(HaveOccurred())
	reader := authorityimpl.NewAdministratorReader(db)
	guard := protect.New(authorityimpl.New(db), reader)
	count, err := reader.Count(ctx, tenant)
	Expect(err).NotTo(HaveOccurred())
	Expect(count).To(BeZero(), "capabilities across different accounts are not an administrator")
	Expect(guard.Run(ctx, tenant, func(context.Context) error { return nil })).To(MatchError(authority.ErrBootstrapRequired))
	Expect(db.Exec("INSERT INTO identity.account_roles(tenant_id, account_id, role_id, assignment_type) VALUES (?, ?, ?, 'SECONDARY'), (?, ?, ?, 'SECONDARY')", tenant, first, assignments, tenant, second, grants).Error).NotTo(HaveOccurred())
	count, err = reader.Count(ctx, tenant)
	Expect(err).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(2)))
	outcomes := make(chan error, 2)
	for _, account := range []string{first, second} {
		go func(accountID string) {
			outcomes <- guard.Run(ctx, tenant, func(tx context.Context) error {
				return transaction.Database(tx, db).Exec("UPDATE identity.accounts SET status = 'SUSPENDED' WHERE tenant_id = ? AND account_id = ?", tenant, accountID).Error
			})
		}(account)
	}
	successes := 0
	for range 2 {
		err := <-outcomes
		if err == nil {
			successes++
		} else {
			Expect(err).To(MatchError(authority.ErrLastAdministrator))
		}
	}
	Expect(successes).To(Equal(1))
	count, err = reader.Count(ctx, tenant)
	Expect(err).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(1)))
	Expect(guard.Run(ctx, tenant, func(tx context.Context) error {
		return transaction.Database(tx, db).Exec("UPDATE identity.tenant_features SET enabled = false WHERE tenant_id = ?", tenant).Error
	})).To(MatchError(authority.ErrLastAdministrator))
	count, err = reader.Count(ctx, tenant)
	Expect(err).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(1)))
	count, err = reader.Count(ctx, "foreign-guard-school")
	Expect(err).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
}
