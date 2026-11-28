package lifecycle_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/role/operations/lifecycle"
)

func TestLifecycle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Role Lifecycle Suite")
}

type guard struct {
	tenant string
	after error
}

func (g *guard) Run(ctx context.Context, tenantID string, operation func(context.Context) error) error {
	g.tenant = tenantID
	if err := operation(ctx); err != nil {
		return err
	}
	return g.after
}

type writer struct {
	request role.Transition
	calls int
}

func (w *writer) Transition(ctx context.Context, request role.Transition) (role.Detail, error) {
	w.calls++
	w.request = request
	return role.Detail{Role: entity.Role{RoleID: request.RoleID, Status: request.Status, Version: request.Version + 1}}, nil
}

func request() role.Transition {
	return role.Transition{TenantID: "school", ActorID: "11111111-1111-4111-8111-111111111111", RoleID: "22222222-2222-4222-8222-222222222222", Version: 1, Status: entity.RoleActive}
}

var _ = Describe("Execute", func() {
	It("runs the transition within the tenant authority guard", func() {
		g := &guard{}
		w := &writer{}
		result, err := lifecycle.New(g, w).Execute(context.Background(), request())
		Expect(err).NotTo(HaveOccurred())
		Expect(g.tenant).To(Equal("school"))
		Expect(w.request).To(Equal(request()))
		Expect(result.Role.Version).To(Equal(int64(2)))
	})
	It("discards the result on a failed authority postcondition", func() {
		result, err := lifecycle.New(&guard{after: authority.ErrLastAdministrator}, &writer{}).Execute(context.Background(), request())
		Expect(err).To(MatchError(authority.ErrLastAdministrator))
		Expect(result).To(Equal(role.Detail{}))
	})
	It("rejects invalid transitions and absent dependencies", func() {
		w := &writer{}
		invalid := request()
		invalid.Status = "UNKNOWN"
		_, err := lifecycle.New(&guard{}, w).Execute(context.Background(), invalid)
		Expect(err).To(MatchError(role.ErrInvalidTransition))
		_, err = lifecycle.New(nil, w).Execute(context.Background(), request())
		Expect(err).To(MatchError(authorization.ErrUnavailable))
		Expect(w.calls).To(BeZero())
	})
})
