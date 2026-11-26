package helpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/approval"
	repo "grovehub/services/identity/internal/app/repository/approval"
	"grovehub/services/identity/internal/app/services/approval/helpers"
)

func TestApprovalHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Approval Helpers Suite")
}

const target = "11111111-1111-4111-8111-111111111111"

var _ = Describe("NormalizeLifecycle", func() {
	It("normalizes the typed payload and produces a stable digest", func() {
		request := dto.CreateLifecycle{Type: kind.AccountLifecycleChange, Intent: dto.LifecycleIntent{TargetID: target, Version: 1, Action: "suspend"}, Reason: " Review ", IdempotencyKey: "key"}
		normalized, code, payload, digest, err := helpers.NormalizeLifecycle(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(code).To(Equal(permission.AccountSuspend))
		Expect(normalized.Reason).To(Equal("Review"))
		Expect(payload).To(MatchJSON(`{"targetId":"11111111-1111-4111-8111-111111111111","version":1,"action":"suspend"}`))
		_, _, _, again, err := helpers.NormalizeLifecycle(normalized)
		Expect(err).NotTo(HaveOccurred())
		Expect(again).To(Equal(digest))
		request.Intent.Version++
		_, _, _, changed, err := helpers.NormalizeLifecycle(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(changed).NotTo(Equal(digest))
	})
	It("rejects unknown owner request types and invalid actions", func() {
		request := dto.CreateLifecycle{Type: kind.PermissionAccess, Intent: dto.LifecycleIntent{TargetID: target, Version: 1, Action: "activate"}, Reason: "Review", IdempotencyKey: "key"}
		_, _, _, _, err := helpers.NormalizeLifecycle(request)
		Expect(err).To(MatchError(repo.ErrInvalid))
		request.Type = kind.RoleLifecycleChange
		request.Intent.Action = "delete"
		_, _, _, _, err = helpers.NormalizeLifecycle(request)
		Expect(err).To(MatchError(repo.ErrInvalid))
	})
})

var _ = Describe("ValidatePolicy", func() {
	It("requires contiguous school-delegable steps", func() {
		policy := repo.Policy{PolicyID: target, Steps: []repo.Step{{Number: 1, Permission: permission.RoleDeactivate}}}
		Expect(helpers.ValidatePolicy(policy)).To(Succeed())
		policy.Steps[0].Number = 2
		Expect(helpers.ValidatePolicy(policy)).To(MatchError(repo.ErrPolicy))
		policy.Steps[0].Number = 1
		policy.Steps[0].Permission = permission.FeatureManage
		Expect(helpers.ValidatePolicy(policy)).To(MatchError(repo.ErrPolicy))
	})
})
