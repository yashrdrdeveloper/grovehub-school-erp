package gateway_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"grovehub/internal/session"
	fixture "grovehub/launcher/internal/resources/tests/helpers/workload"
	mock "grovehub/launcher/internal/resources/tests/mocks/gateway/external"
	external "grovehub/launcher/internal/src/services/gateway/operations/external/impl"
)

// @Testing suite=LauncherExternalGateway
func SecuritySpecs() {
	Describe("External gateway security", func() {
		It("serves only registered health probes without workload verification", func() {
			mux := http.NewServeMux()
			mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(204)
			})
			mux.HandleFunc("GET /ready", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(503)
			})
			gateway := external.New(mux, nil, nil, "")
			for _, tc := range []struct {
				path   string
				status int
			}{{"/health", 204}, {"/ready", 503}} {
				w := httptest.NewRecorder()
				gateway.ServeHTTP(w, httptest.NewRequest("GET", tc.path, nil))
				Expect(w.Code).To(Equal(tc.status))
			}
		})
		It("rejects missing and invalid browser credentials", func() {
			mux := http.NewServeMux()
			mux.HandleFunc("GET /foundation/protected/persons", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(204)
			})
			access := mock.NewMockAccessVerifier(gomock.NewController(GinkgoT()))
			access.EXPECT().VerifyAccess(gomock.Any(), "bad").Return(session.Principal{}, errors.New("invalid"))
			access.EXPECT().VerifyAccess(gomock.Any(), "anonymous").Return(session.Principal{}, nil)
			for _, tc := range []struct {
				secret     string
				noVerifier bool
			}{{"", true}, {"", false}, {"bad", false}, {"anonymous", false}} {
				request := httptest.NewRequest("GET", "/foundation/protected/persons", nil)
				if tc.secret != "" {
					request.AddCookie(&http.Cookie{Name: session.AccessCookie, Value: tc.secret})
				}

				w := httptest.NewRecorder()
				if tc.noVerifier {
					external.New(mux, nil, nil, "").ServeHTTP(w, request)
				} else {
					external.New(mux, access, nil, "").ServeHTTP(w, request)
				}

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			}
		})
		It("requires CSRF for public refresh and protected mutation routes", func() {
			const origin = "https://school.example"
			mux := http.NewServeMux()
			mux.HandleFunc("POST /identity/public/auth/refresh", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(204)
			})
			mux.HandleFunc("POST /foundation/protected/persons", func(w http.ResponseWriter, r *http.Request) {
				principal, ok := session.PrincipalFrom(r.Context())
				Expect(ok).To(BeTrue())
				Expect(principal.AccountID).To(Equal("account"))
				w.WriteHeader(204)
			})
			access := mock.NewMockAccessVerifier(gomock.NewController(GinkgoT()))
			access.EXPECT().VerifyAccess(gomock.Any(), "secret").Return(session.Principal{AccountID: "account", SessionID: "session"}, nil).Times(2)
			gateway := external.New(mux, access, nil, origin)
			for _, path := range []string{"/identity/public/auth/refresh", "/foundation/protected/persons"} {
				for _, allowed := range []bool{false, true} {
					request := httptest.NewRequest("POST", path, nil)
					request.Header.Set("Origin", origin)
					request.AddCookie(&http.Cookie{Name: session.AccessCookie, Value: "secret"})
					if allowed {
						request.AddCookie(&http.Cookie{Name: session.CSRFName, Value: "token"})
						request.Header.Set("X-CSRF-Token", "token")
					}

					w := httptest.NewRecorder()
					gateway.ServeHTTP(w, request)
					if allowed {
						Expect(w.Code).To(Equal(204))
					} else {
						Expect(w.Code).To(Equal(403))
					}
				}
			}
		})
		It("requires a verified workload context and adapts only valid identity proofs", func() {
			const pattern = "POST /identity/restricted/auth/refresh"
			mux := http.NewServeMux()
			mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
				cookie, err := r.Cookie(session.RefreshCookie)
				Expect(err).NotTo(HaveOccurred())
				Expect(cookie.Value).To(Equal("proof"))
				w.WriteHeader(204)
			})
			verified, err := fixture.VerifiedContext("school")
			Expect(err).NotTo(HaveOccurred())
			workload := mock.NewMockWorkloadVerifier(gomock.NewController(GinkgoT()))
			gomock.InOrder(
				workload.EXPECT().VerifyWorkload(gomock.Any(), gomock.Any()).Return(nil, errors.New("denied")),
				workload.EXPECT().VerifyWorkload(gomock.Any(), gomock.Any()).Return(nil, nil),
				workload.EXPECT().VerifyWorkload(gomock.Any(), gomock.Any()).Return(context.Background(), nil),
				workload.EXPECT().VerifyWorkload(gomock.Any(), gomock.Any()).Return(verified, nil),
			)
			gateway := external.New(mux, nil, workload, "")
			for i := 0; i < 4; i++ {
				r := httptest.NewRequest("POST", "/identity/restricted/auth/refresh", nil)
				r.Header.Set("X-GH-Identity-Proof", "proof")
				w := httptest.NewRecorder()
				gateway.ServeHTTP(w, r)
				if i == 3 {
					Expect(w.Code).To(Equal(204))
				} else {
					Expect(w.Code).To(Equal(403))
				}
			}
		})
	})
}
