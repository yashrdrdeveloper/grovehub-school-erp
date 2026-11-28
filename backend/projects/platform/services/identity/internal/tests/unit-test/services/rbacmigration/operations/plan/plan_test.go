package plan_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/internal/app/constants/persona"
	"grovehub/services/identity/internal/app/services/rbacmigration"
	"grovehub/services/identity/internal/app/services/rbacmigration/operations/plan"
)

func TestMigrationPlan(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Legacy RBAC Migration Plan Suite")
}

const accountID = "11111111-1111-4111-8111-111111111111"

var _ = Describe("RoleID", func() {
	It("is deterministic and tenant and persona specific", func() {
		id := plan.RoleID("school-a", persona.Teacher)
		Expect(id).To(Equal(plan.RoleID("school-a", persona.Teacher)))
		Expect(id).NotTo(Equal(plan.RoleID("school-b", persona.Teacher)))
		Expect(id).NotTo(Equal(plan.RoleID("school-a", persona.Accountant)))
	})
})

var _ = Describe("Build", func() {
	It("maps legacy codes while preserving versions and reporting unassigned accounts", func() {
		primary := string(persona.Teacher)
		inventory := rbacmigration.Inventory{TenantID: "school", Accounts: []rbacmigration.LegacyAccount{
			{AccountID: accountID, Version: 3, Primary: &primary, Secondary: []string{string(persona.Coordinator)}},
			{AccountID: "22222222-2222-4222-8222-222222222222", Version: 1},
		}}
		result, err := plan.Build(inventory)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Roles).To(HaveLen(6))
		Expect(result.Assignments).To(Equal([]rbacmigration.AssignmentPlan{{AccountID: accountID, Version: 3, PrimaryRoleID: plan.RoleID("school", persona.Teacher), SecondaryRoleIDs: []string{plan.RoleID("school", persona.Coordinator)}}}))
		Expect(result.UnassignedAccountIDs).To(HaveLen(1))
		again, err := plan.Build(inventory)
		Expect(err).NotTo(HaveOccurred())
		Expect(again).To(Equal(result))
	})
	It("rejects unknown and duplicate legacy roles", func() {
		primary := string(persona.Teacher)
		inventory := rbacmigration.Inventory{TenantID: "school", Accounts: []rbacmigration.LegacyAccount{{AccountID: accountID, Version: 1, Primary: &primary, Secondary: []string{primary}}}}
		_, err := plan.Build(inventory)
		Expect(err).To(MatchError(rbacmigration.ErrInvalidInventory))
		inventory.Accounts[0].Secondary = []string{"UNKNOWN"}
		_, err = plan.Build(inventory)
		Expect(err).To(MatchError(rbacmigration.ErrInvalidInventory))
	})
	It("refuses to overwrite normalized assignments or adopt existing named roles", func() {
		inventory := rbacmigration.Inventory{TenantID: "school", NormalizedAssignments: 1}
		_, err := plan.Build(inventory)
		Expect(err).To(MatchError(rbacmigration.ErrConflict))
		inventory.NormalizedAssignments = 0
		inventory.Roles = []rbacmigration.ExistingRole{{NormalizedName: "teacher"}}
		_, err = plan.Build(inventory)
		Expect(err).To(MatchError(rbacmigration.ErrConflict))
	})
})
