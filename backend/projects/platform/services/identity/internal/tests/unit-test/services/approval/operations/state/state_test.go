package state_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/services/approval/operations/state"
)

func TestApprovalState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Approval State Suite")
}

var now = time.Date(2026, 11, 26, 12, 0, 0, 0, time.UTC)

func pending() state.Snapshot {
	return state.Snapshot{TenantID: "school", RequesterID: "requester", Status: approval.Pending, CurrentStep: 1, StepCount: 2, Version: 1, ExpiresAt: now.Add(time.Hour)}
}

var reviewer = state.Actor{TenantID: "school", AccountID: "reviewer"}

var _ = Describe("Decide", func() {
	It("progresses sequentially without applying the requested effect", func() {
		first, err := state.Decide(pending(), reviewer, 1, approval.Approve, "", now)
		Expect(err).NotTo(HaveOccurred())
		Expect(first.Status).To(Equal(approval.Pending))
		Expect(first.CurrentStep).To(Equal(2))
		Expect(first.Version).To(Equal(int64(2)))
		last, err := state.Decide(first, reviewer, 2, approval.Approve, "", now)
		Expect(err).NotTo(HaveOccurred())
		Expect(last.Status).To(Equal(approval.Approved))
		Expect(last.Version).To(Equal(int64(3)))
		_, err = state.Decide(last, reviewer, 3, approval.Approve, "", now)
		Expect(err).To(MatchError(state.ErrInvalid))
	})
	It("rejects stale versions, self-approval, foreign reviewers and expired requests", func() {
		current := pending()
		_, err := state.Decide(current, reviewer, 2, approval.Approve, "", now)
		Expect(err).To(MatchError(state.ErrVersion))
		_, err = state.Decide(current, state.Actor{TenantID: "school", AccountID: "requester"}, 1, approval.Approve, "", now)
		Expect(err).To(MatchError(state.ErrForbidden))
		_, err = state.Decide(current, state.Actor{TenantID: "foreign", AccountID: "reviewer"}, 1, approval.Approve, "", now)
		Expect(err).To(MatchError(state.ErrForbidden))
		_, err = state.Decide(current, reviewer, 1, approval.Approve, "", current.ExpiresAt)
		Expect(err).To(MatchError(state.ErrExpired))
	})
	It("requires a rejection reason", func() {
		_, err := state.Decide(pending(), reviewer, 1, approval.Reject, " ", now)
		Expect(err).To(MatchError(state.ErrInvalid))
		result, err := state.Decide(pending(), reviewer, 1, approval.Reject, "Not authorized", now)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Status).To(Equal(approval.Rejected))
	})
})

var _ = Describe("Cancel", func() {
	It("allows only the requester to cancel a pending request", func() {
		_, err := state.Cancel(pending(), reviewer, 1, now)
		Expect(err).To(MatchError(state.ErrForbidden))
		result, err := state.Cancel(pending(), state.Actor{TenantID: "school", AccountID: "requester"}, 1, now)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Status).To(Equal(approval.Cancelled))
	})
})

var _ = Describe("Expire", func() {
	It("expires pending requests at the deadline", func() {
		current := pending()
		_, err := state.Expire(current, 1, now)
		Expect(err).To(MatchError(state.ErrInvalid))
		result, err := state.Expire(current, 1, current.ExpiresAt)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Status).To(Equal(approval.Expired))
	})
})

var _ = Describe("EffectResult", func() {
	It("separates approval from applied, failed and stale outcomes", func() {
		current := pending()
		_, err := state.EffectResult(current, 1, approval.Applied, now)
		Expect(err).To(MatchError(state.ErrInvalid))
		current.Status = approval.Approved
		current.CurrentStep = current.StepCount
		for _, outcome := range []approval.Status{approval.Applied, approval.Failed, approval.Stale} {
			result, err := state.EffectResult(current, 1, outcome, now)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Status).To(Equal(outcome))
			Expect(result.Version).To(Equal(int64(2)))
		}
		_, err = state.EffectResult(current, 1, approval.Applied, current.ExpiresAt)
		Expect(err).To(MatchError(state.ErrExpired))
	})
})
