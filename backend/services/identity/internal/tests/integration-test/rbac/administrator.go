package rbac_test

import (
	"context"

	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/repository/authority"
	assignmentreplace "grovehub/services/identity/internal/app/services/accountrole/operations/replace"
	accountrolerepo "grovehub/services/identity/internal/app/repository/accountrole"
	accountroleimpl "grovehub/services/identity/internal/app/repository/accountrole/impl"
	"grovehub/services/identity/internal/app/constants/persona"
	"grovehub/services/identity/internal/app/entity"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	roleimpl "grovehub/services/identity/internal/app/repository/role/impl"
	"grovehub/services/identity/internal/app/services/role/operations/lifecycle"
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
	roles := roleimpl.New(db)
	draft, err := roles.CreateDraft(ctx, rolerepo.Draft{TenantID: tenant, ActorID: first, Name: "Lifecycle fixture", Personas: []persona.Code{persona.Teacher}})
	Expect(err).NotTo(HaveOccurred())
	operation := lifecycle.New(guard, roleimpl.NewLifecycleWriter(db))
	transition := rolerepo.Transition{TenantID: tenant, ActorID: first, RoleID: draft.Role.RoleID, Version: 1, Status: entity.RoleActive}
	_, err = operation.Execute(ctx, transition)
	Expect(err).To(MatchError(rolerepo.ErrNotReady))
	_, err = roleimpl.NewDraftWriter(db).ReplaceDraftPermissions(ctx, rolerepo.DraftPermissionUpdate{TenantID: tenant, ActorID: first, RoleID: draft.Role.RoleID, Version: 1, Codes: []permission.Code{permission.RolePermissionManage}})
	Expect(err).NotTo(HaveOccurred())
	transition.Version = 2
	active, err := operation.Execute(ctx, transition)
	Expect(err).NotTo(HaveOccurred())
	Expect(active.Role.Status).To(Equal(entity.RoleActive))
	Expect(active.Role.Version).To(Equal(int64(3)))
	_, err = operation.Execute(ctx, transition)
	Expect(err).To(MatchError(rolerepo.ErrVersionConflict))
	transition.Version = 3
	_, err = operation.Execute(ctx, transition)
	Expect(err).To(MatchError(rolerepo.ErrInvalidTransition))
	transition.Status = entity.RoleInactive
	Expect(db.Exec("INSERT INTO identity.account_roles(tenant_id, account_id, role_id, assignment_type) VALUES (?, ?, ?, 'SECONDARY')", tenant, first, transition.RoleID).Error).NotTo(HaveOccurred())
	_, err = operation.Execute(ctx, transition)
	Expect(err).To(MatchError(rolerepo.ErrAssigned))
	Expect(db.Exec("DELETE FROM identity.account_roles WHERE tenant_id = ? AND role_id = ?", tenant, transition.RoleID).Error).NotTo(HaveOccurred())
	inactive, err := operation.Execute(ctx, transition)
	Expect(err).NotTo(HaveOccurred())
	Expect(inactive.Role.Status).To(Equal(entity.RoleInactive))
	Expect(inactive.Role.Version).To(Equal(int64(4)))
	transition.RoleID = grants
	transition.Version = 1
	_, err = operation.Execute(ctx, transition)
	Expect(err).To(MatchError(rolerepo.ErrAssigned))
	assignmentsRepo := accountroleimpl.New(db)
	assignment, err := assignmentsRepo.Find(ctx, tenant, first)
	Expect(err).NotTo(HaveOccurred())
	Expect(assignment.PrimaryRoleID).To(Equal(grants))
	Expect(assignment.SecondaryRoleIDs).To(Equal([]string{assignments}))
	_, err = assignmentsRepo.Find(ctx, "foreign-guard-school", first)
	Expect(err).To(MatchError(accountrolerepo.ErrNotFound))
	candidates, err := assignmentsRepo.Candidates(ctx, tenant, []string{grants, assignments})
	Expect(err).NotTo(HaveOccurred())
	Expect(candidates).To(HaveLen(2))
	_, err = assignmentsRepo.Candidates(ctx, "foreign-guard-school", []string{grants})
	Expect(err).To(MatchError(accountrolerepo.ErrNotFound))
	// Select the sole active administrator after the concurrent suspension test.
	var survivor string
	Expect(db.Raw("SELECT account_id::text FROM identity.accounts WHERE tenant_id = ? AND status = 'ACTIVE'", tenant).Scan(&survivor).Error).NotTo(HaveOccurred())
	before, err := assignmentsRepo.Find(ctx, tenant, survivor)
	Expect(err).NotTo(HaveOccurred())
	replacement := accountrolerepo.Replacement{
		TenantID: tenant, ActorID: survivor, AccountID: survivor, Version: before.Version,
		PrimaryRoleID: before.PrimaryRoleID, SecondaryRoleIDs: []string{},
		RoleVersions: []accountrolerepo.RoleVersion{{RoleID: grants, Version: 1}, {RoleID: assignments, Version: 1}},
	}
	replaceOperation := assignmentreplace.New(guard, accountroleimpl.NewWriter(db))
	_, err = replaceOperation.Execute(ctx, replacement)
	Expect(err).To(MatchError(authority.ErrLastAdministrator))
	after, err := assignmentsRepo.Find(ctx, tenant, survivor)
	Expect(err).NotTo(HaveOccurred())
	Expect(after).To(Equal(before))
	replacement.PrimaryRoleID = before.SecondaryRoleIDs[0]
	replacement.SecondaryRoleIDs = []string{before.PrimaryRoleID}
	after, err = replaceOperation.Execute(ctx, replacement)
	Expect(err).NotTo(HaveOccurred())
	Expect(after.Version).To(Equal(before.Version + 1))
	Expect(after.PrimaryRoleID).To(Equal(replacement.PrimaryRoleID))
	Expect(after.SecondaryRoleIDs).To(Equal(replacement.SecondaryRoleIDs))
	_, err = replaceOperation.Execute(ctx, replacement)
	Expect(err).To(MatchError(accountrolerepo.ErrVersionConflict))
	replacement.Version = after.Version
	unchanged, err := replaceOperation.Execute(ctx, replacement)
	Expect(err).NotTo(HaveOccurred())
	Expect(unchanged.Version).To(Equal(after.Version))
	replacement.RoleVersions[0].Version = 99
	_, err = replaceOperation.Execute(ctx, replacement)
	Expect(err).To(MatchError(accountrolerepo.ErrVersionConflict))
	count, err = reader.Count(ctx, "foreign-guard-school")
	Expect(err).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
}
