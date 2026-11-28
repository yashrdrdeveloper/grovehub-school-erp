package session_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "grovehub/libs/platform/session"
)

// @Testing suite=LauncherSession
func SessionSpecs() {
	Describe("Tenant generation", func() {
		It("rejects disabled and stale sessions without accepting a different tenant", func() {
			generation := int64(1)
			active := true
			ctx := WithTenantCheck(context.Background(), func(_ context.Context, tenant string) (int64, error) {
				if tenant != "school" || !active {
					return 0, ErrTenantDenied
				}
				return generation, nil
			})
			issued, err := TenantGeneration(ctx, "school")
			Expect(err).NotTo(HaveOccurred())
			Expect(issued).To(Equal(int64(1)))
			Expect(CheckTenantSession(ctx, "school", issued)).To(Succeed())
			active = false
			generation++
			Expect(CheckTenantSession(ctx, "school", issued)).To(MatchError(ErrTenantDenied))
			active = true
			for _, old := range []int64{0, issued, generation + 1} {
				Expect(CheckTenantSession(ctx, "school", old)).To(MatchError(ErrTenantDenied))
			}

			Expect(CheckTenantSession(ctx, "school", generation)).To(Succeed())
			Expect(CheckTenantSession(ctx, "other", generation)).To(MatchError(ErrTenantDenied))
		})
		It("derives tenant from stored account identity", func() {
			tenant := "stored-school"
			ctx := WithTenantCheck(context.Background(), func(_ context.Context, got string) (int64, error) {
				Expect(got).To(Equal(tenant))
				return 0, ErrTenantDenied
			})
			Expect(CheckAccountTenant(ctx, &tenant)).To(MatchError(ErrTenantDenied))
			Expect(CheckAccountTenant(ctx, nil)).To(MatchError(ErrTenantDenied))
		})
		It("handles missing and inactive policies, lookup errors and valid generations", func() {
			ctx := context.Background()
			Expect(ActiveTenantCheck(nil)).To(BeNil())
			Expect(WithTenantCheck(ctx, nil)).To(Equal(ctx))
			n, err := TenantGeneration(ctx, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(int64(0)))
			Expect(CheckTenantSession(ctx, "school", 3)).To(Succeed())
			check := ActiveTenantCheck(func(_ context.Context, id string) (TenantState, error) {
				switch id {
				case "active":
					return TenantState{Active: true, Generation: 3}, nil
				case "disabled":
					return TenantState{Generation: 3}, nil
				case "zero":
					return TenantState{Active: true}, nil
				default:
					return TenantState{}, errors.New("unavailable")
				}
			})
			scoped := WithTenantCheck(ctx, check)
			for _, id := range []string{"", "disabled", "zero", "missing"} {
				Expect(CheckTenant(scoped, id)).To(MatchError(ErrTenantDenied))
			}

			id := "active"
			Expect(CheckAccountTenant(scoped, &id)).To(Succeed())
			Expect(CheckAccountGeneration(scoped, &id, 3)).To(Succeed())
			Expect(CheckAccountGeneration(scoped, nil, 3)).To(MatchError(ErrTenantDenied))
		})
	})
	Describe("Principal", func() {
		It("accepts only principals with both account and session", func() {
			ctx := context.Background()
			_, ok := PrincipalFrom(ctx)
			Expect(ok).To(BeFalse())
			_, ok = PrincipalFrom(WithPrincipal(ctx, Principal{AccountID: "a"}))
			Expect(ok).To(BeFalse())
			_, ok = PrincipalFrom(WithPrincipal(ctx, Principal{SessionID: "s"}))
			Expect(ok).To(BeFalse())
			p := Principal{AccountID: "a", SessionID: "s", TenantID: "school", SessionGeneration: 1}
			got, ok := PrincipalFrom(WithPrincipal(ctx, p))
			Expect(ok).To(BeTrue())
			Expect(got).To(Equal(p))
		})
	})
}
