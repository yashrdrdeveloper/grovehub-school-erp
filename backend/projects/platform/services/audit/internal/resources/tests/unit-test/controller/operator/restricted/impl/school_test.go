package impl_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	"grovehub/services/audit/internal/app/controller/operator/constants"
	"grovehub/services/audit/internal/app/controller/operator/restricted"
	controller "grovehub/services/audit/internal/app/controller/operator/restricted/impl"
	"grovehub/services/audit/internal/app/services/operator"
)

type credentialVerifier struct {
	claims gateway.WorkloadClaims
}

func (v credentialVerifier) VerifyCredential(_ context.Context, _ string) (gateway.WorkloadClaims, error) {
	return v.claims, nil
}

type schoolService struct {
	operator.Service
	exists bool
	err    error
	calls  int
	school string
}

func (s *schoolService) SchoolCheck(_ context.Context, school string) (bool, error) {
	s.calls++
	s.school = school
	return s.exists, s.err
}

func TestSchoolCheck(t *testing.T) {
	const school = "e93696cc-15f0-4e06-9577-0d5cfc2a83c8"
	for _, tc := range []struct {
		name    string
		id      string
		subject string
		scopes  []string
		exists  bool
		err     error
		status  int
		calls   int
	}{
		{"exists", school, constants.SchoolSubject + school, []string{constants.SchoolCheckScope}, true, nil, http.StatusNoContent, 1},
		{"absent", school, constants.SchoolSubject + school, []string{constants.SchoolCheckScope}, false, nil, http.StatusNotFound, 1},
		{"unavailable", school, constants.SchoolSubject + school, []string{constants.SchoolCheckScope}, false, errors.New("database failure"), http.StatusServiceUnavailable, 1},
		{"cross-tenant", school, constants.SchoolSubject + "other", []string{constants.SchoolCheckScope}, true, nil, http.StatusForbidden, 0},
		{"invalid UUID", "invalid", constants.SchoolSubject + "invalid", []string{constants.SchoolCheckScope}, true, nil, http.StatusForbidden, 0},
		{"extra scope", school, constants.SchoolSubject + school, []string{constants.SchoolCheckScope, "other"}, true, nil, http.StatusForbidden, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := &schoolService{exists: tc.exists, err: tc.err}
			mux := http.NewServeMux()
			restricted.Register(mux, controller.New(s))

			r := httptest.NewRequest(http.MethodGet, "/audit/restricted/schools/"+tc.id+"/activity-check", nil)
			r.Pattern = constants.SchoolCheckPattern
			r.Header.Set("Authorization", "Bearer fixture")
			verifier := gateway.ScopedWorkloadVerifier{
				Credentials: credentialVerifier{claims: gateway.WorkloadClaims{
					Subject: tc.subject, Audience: "platform", Scopes: tc.scopes,
					ExpiresAt: time.Now().Add(time.Minute),
				}},
				Audience: "platform",
				Routes: map[string]gateway.RoutePolicy{
					constants.SchoolCheckPattern: {Scope: constants.SchoolCheckScope},
				},
			}
			ctx, err := verifier.VerifyWorkload(r.Context(), r)
			if err != nil {
				t.Fatalf("fixture verification failed: %v", err)
			}

			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r.WithContext(ctx))
			if w.Code != tc.status || s.calls != tc.calls || w.Body.Len() != 0 {
				t.Fatalf("status=%d calls=%d body=%q", w.Code, s.calls, w.Body.String())
			}
			if w.Header().Get("Cache-Control") != "no-store" {
				t.Fatal("missing cache protection")
			}
			if s.calls > 0 && s.school != tc.id {
				t.Fatalf("wrong school passed to service: %s", s.school)
			}
		})
	}

	s := &schoolService{}
	mux := http.NewServeMux()
	restricted.Register(mux, controller.New(s))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/audit/restricted/schools/"+school+"/activity-check", nil))
	if w.Code != http.StatusForbidden || s.calls != 0 {
		t.Fatal("unauthenticated request reached the service")
	}
}
