package workload_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"time"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/workload"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type credentialFunc func(context.Context, string) (contracts.WorkloadClaims, error)

func (f credentialFunc) VerifyCredential(ctx context.Context, token string) (contracts.WorkloadClaims, error) {
	return f(ctx, token)
}

var _ = Describe("VerifyWorkload", func() {
	const route = "GET /service/restricted/read"
	var claims contracts.WorkloadClaims
	var verifier workload.ScopedWorkloadVerifier

	BeforeEach(func() {
		claims = contracts.WorkloadClaims{Subject:"service", Audience:"platform", ExpiresAt:time.Now().Add(time.Minute), Scopes:[]string{"read"}, TenantID:"tenant", ActorGID:"actor"}
		verifier = workload.ScopedWorkloadVerifier{
			Audience:"platform", Routes:map[string]contracts.RoutePolicy{route:{Scope:"read", RequireActor:true}},
			Credentials:credentialFunc(func(_ context.Context, token string)(contracts.WorkloadClaims,error){
				if token != "valid" { return contracts.WorkloadClaims{}, errors.New("invalid") }
				return claims,nil
			}),
		}
	})

	It("stores verified claims and isolates scope slices from callers", func() {
		request := httptest.NewRequest("GET", "/service/restricted/read", nil)
		request.Pattern = route
		request.Header.Set("Authorization","Bearer valid")
		ctx, err := verifier.VerifyWorkload(context.Background(), request)
		Expect(err).NotTo(HaveOccurred())
		claims.Scopes[0] = "write"
		got, ok := workload.WorkloadFrom(ctx)
		Expect(ok).To(BeTrue())
		Expect(got.Scopes).To(Equal([]string{"read"}))
		got.Scopes[0] = "other"
		again, ok := workload.WorkloadFrom(ctx)
		Expect(ok).To(BeTrue())
		Expect(again.Scopes).To(Equal([]string{"read"}))
		_, ok = workload.WorkloadFrom(nil)
		Expect(ok).To(BeFalse())
		_, ok = workload.WorkloadFrom(context.Background())
		Expect(ok).To(BeFalse())
	})

	It("rejects missing dependencies, malformed headers, unknown routes and invalid verified claims", func() {
		request := httptest.NewRequest("GET", "/service/restricted/read", nil)
		request.Pattern = route
		request.Header.Set("Authorization","Bearer valid")
		for _, bad := range []workload.ScopedWorkloadVerifier{{}, {Audience:"platform", Credentials:verifier.Credentials}} {
			_, err := bad.VerifyWorkload(context.Background(), request)
			Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		}
		_, err := verifier.VerifyWorkload(nil, request)
		Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		_, err = verifier.VerifyWorkload(context.Background(), nil)
		Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		for _, value := range []string{"", "Basic valid", "Bearer", "Bearer ", "Bearer bad token", "Bearer invalid"} {
			request.Header.Set("Authorization", value)
			_, err = verifier.VerifyWorkload(context.Background(), request)
			Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		}
		request.Header.Set("Authorization","Bearer valid")
		request.Header.Add("Authorization","Bearer valid")
		_, err = verifier.VerifyWorkload(context.Background(), request)
		Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		request.Header.Set("Authorization","Bearer valid")
		original := claims
		for _, mutate := range []func(*contracts.WorkloadClaims){
			func(v *contracts.WorkloadClaims){v.Subject=""}, func(v *contracts.WorkloadClaims){v.Audience="other"},
			func(v *contracts.WorkloadClaims){v.ExpiresAt=time.Now().Add(-time.Minute)}, func(v *contracts.WorkloadClaims){v.TenantID=""},
			func(v *contracts.WorkloadClaims){v.TenantID="";v.ActorGID=""}, func(v *contracts.WorkloadClaims){v.Scopes=[]string{"write"}},
		} {
			claims = original
			mutate(&claims)
			_, err = verifier.VerifyWorkload(context.Background(), request)
			Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		}
	})
})
