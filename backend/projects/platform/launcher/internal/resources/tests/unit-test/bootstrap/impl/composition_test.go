package impl_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	mock "grovehub/launcher/internal/resources/tests/mocks/bootstrap"
	"grovehub/launcher/internal/src/bootstrap"
	"grovehub/launcher/internal/src/bootstrap/impl"
	"grovehub/launcher/internal/src/config"
	"grovehub/launcher/internal/src/server"
	"grovehub/libs/platform/session"
	"grovehub/services/audit"
	"grovehub/services/foundation"
	person "grovehub/services/foundation/client/person"
	"grovehub/services/identity"
	"grovehub/services/infrastructure"
	"grovehub/services/infrastructure/client/messaging"
	"grovehub/services/integration"
	connection "grovehub/services/integration/client/connection"
)

type compositionFixture struct {
	dependencies            bootstrap.Dependencies
	businessMock, auditMock sqlmock.Sqlmock
	audit                   *mock.MockAuditModule
	messaging               *mock.MockMessagingClient
	server                  *mock.MockServer
	handler                 http.Handler
	draining                *atomic.Bool
	started, closed         []string
}

func mockPool() (*gorm.DB, sqlmock.Sqlmock) {
	raw, expect, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	Expect(err).NotTo(HaveOccurred())
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: raw}), &gorm.Config{DisableAutomaticPing: true})
	Expect(err).NotTo(HaveOccurred())
	DeferCleanup(func() {
		_ = raw.Close()
	})
	return db, expect
}

func newCompositionFixture(failure string) *compositionFixture {
	ctrl := gomock.NewController(GinkgoT())
	f := &compositionFixture{dependencies: impl.DefaultDependencies()}
	db, businessMock := mockPool()
	auditDB, auditMock := mockPool()
	f.businessMock, f.auditMock = businessMock, auditMock
	infra := mock.NewMockInfrastructureModule(ctrl)
	foundationModule := mock.NewMockFoundationModule(ctrl)
	integrationModule := mock.NewMockIntegrationModule(ctrl)
	identityModule := mock.NewMockIdentityModule(ctrl)
	f.audit = mock.NewMockAuditModule(ctrl)
	f.messaging = mock.NewMockMessagingClient(ctrl)
	f.server = mock.NewMockServer(ctrl)

	record := func(name string) error {
		f.started = append(f.started, name)
		if failure == name {
			return errors.New("injected " + name)
		}

		return nil
	}
	closeModule := func(name string) {
		f.closed = append(f.closed, name)
	}
	f.dependencies.OpenDatabase = func(context.Context) (*gorm.DB, error) {
		if err := record("database"); err != nil {
			return nil, err
		}

		if failure == "database pool" {
			return &gorm.DB{Config: &gorm.Config{}}, nil
		}

		return db, nil
	}
	f.dependencies.NewInfrastructure = func(_ messaging.Config, got *gorm.DB, _ *slog.Logger) (infrastructure.Module, error) {
		Expect(got).To(BeIdenticalTo(db))
		if err := record("infrastructure"); err != nil {
			return nil, err
		}

		return infra, nil
	}
	f.dependencies.NewFoundation = func(_ *slog.Logger, got *gorm.DB, _ messaging.EventStore) (foundation.Module, error) {
		Expect(got).To(BeIdenticalTo(db))
		if err := record("foundation"); err != nil {
			return nil, err
		}

		return foundationModule, nil
	}
	f.dependencies.NewIntegration = func(_ *slog.Logger, got *gorm.DB, events messaging.MessagingClient, _ session.TenantStateLookup) (integration.Module, error) {
		Expect(got).To(BeIdenticalTo(db))
		Expect(events).To(BeIdenticalTo(f.messaging))
		if err := record("integration"); err != nil {
			return nil, err
		}

		return integrationModule, nil
	}
	f.dependencies.NewIdentity = func(_ *slog.Logger, got *gorm.DB, _ person.Client, _ connection.Client, _ messaging.EventStore, _ ...session.TenantStateLookup) (identity.Module, error) {
		Expect(got).To(BeIdenticalTo(db))
		if err := record("identity"); err != nil {
			return nil, err
		}

		return identityModule, nil
	}
	f.dependencies.OpenAuditDatabase = func(_ context.Context, shared *gorm.DB) (*gorm.DB, func(), error) {
		Expect(shared).To(BeIdenticalTo(db))
		if err := record("audit database"); err != nil {
			return nil, nil, err
		}

		if failure == "audit pool" {
			return &gorm.DB{Config: &gorm.Config{}}, func() {
				closeModule("audit database")
			}, nil
		}

		return auditDB, func() {
			closeModule("audit database")
		}, nil
	}
	f.dependencies.NewAudit = func(_ *slog.Logger, _ *gorm.DB, events messaging.MessagingClient, access session.ProtectedRouteAuthorizer) (audit.Module, error) {
		Expect(events).To(BeIdenticalTo(f.messaging))
		Expect(access).To(BeIdenticalTo(identityModule))
		if err := record("audit"); err != nil {
			return nil, err
		}

		return f.audit, nil
	}
	f.dependencies.StartTenantReconciler = func(_ context.Context, _ *slog.Logger, _ *gorm.DB, lookup session.TenantStateLookup) func() {
		Expect(lookup).NotTo(BeNil())
		f.started = append(f.started, "tenant worker")
		return func() {
			closeModule("tenant worker")
		}
	}
	f.dependencies.NewServer = func(_ string, handler http.Handler, _ *slog.Logger, draining *atomic.Bool) server.Server {
		f.started = append(f.started, "server")
		f.handler, f.draining = handler, draining
		return f.server
	}
	infra.EXPECT().EventStore().Return(nil).AnyTimes()
	infra.EXPECT().MessagingClient().Return(f.messaging).AnyTimes()
	foundationModule.EXPECT().Client().Return(nil).AnyTimes()
	integrationModule.EXPECT().Client().Return(nil).AnyTimes()
	infra.EXPECT().Close().DoAndReturn(func() error {
		closeModule("infrastructure")
		return nil
	}).AnyTimes()
	foundationModule.EXPECT().Close().Do(func() {
		closeModule("foundation")
	}).AnyTimes()
	integrationModule.EXPECT().Close().Do(func() {
		closeModule("integration")
	}).AnyTimes()
	identityModule.EXPECT().Close().Do(func() {
		closeModule("identity")
	}).AnyTimes()
	f.audit.EXPECT().Close().Do(func() {
		closeModule("audit")
	}).AnyTimes()
	foundationModule.EXPECT().RegisterHTTP(gomock.Any()).AnyTimes()
	integrationModule.EXPECT().RegisterHTTP(gomock.Any()).AnyTimes()
	identityModule.EXPECT().RegisterHTTP(gomock.Any()).AnyTimes()
	f.audit.EXPECT().RegisterHTTP(gomock.Any()).AnyTimes()
	identityModule.EXPECT().RouteLimits().Return(nil).AnyTimes()
	return f
}

func compositionLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// @Testing suite=LauncherBootstrapImplementation
func CompositionSpecs() {
	Describe("Platform composition", func() {
		It("stops at each failed dependency and rolls back successfully opened modules in reverse order", func() {
			for _, tc := range []struct {
				failure string
				closed  []string
			}{
				{"database", nil},
				{"database pool", nil},
				{"infrastructure", nil},
				{"foundation", []string{"infrastructure"}},
				{"integration", []string{"foundation", "infrastructure"}},
				{"identity", []string{"integration", "foundation", "infrastructure"}},
				{"audit database", []string{"identity", "integration", "foundation", "infrastructure"}},
				{"audit", []string{"audit database", "identity", "integration", "foundation", "infrastructure"}},
				{"audit pool", []string{"audit", "audit database", "identity", "integration", "foundation", "infrastructure"}},
				{"gateway", []string{"audit", "audit database", "identity", "integration", "foundation", "infrastructure"}},
			} {
				By(tc.failure)
				fixture := newCompositionFixture(tc.failure)
				cfg := config.Config{Address: "127.0.0.1:0"}
				if tc.failure == "gateway" {
					cfg.Workload = "{"
				}

				runtime, err := impl.NewWithDependencies(compositionLogger(), cfg, fixture.dependencies).Build()
				Expect(err).To(HaveOccurred())
				Expect(runtime).To(BeNil())
				Expect(fixture.closed).To(Equal(tc.closed))
			}
		})
		It("builds all modules, checks readiness of every dependency and drains before closing", func() {
			fixture := newCompositionFixture("")
			cfg := config.Config{Address: "127.0.0.1:0"}
			runtime, err := impl.NewWithDependencies(compositionLogger(), cfg, fixture.dependencies).Build()
			Expect(err).NotTo(HaveOccurred())
			Expect(fixture.started).To(Equal([]string{"database", "infrastructure", "foundation", "integration", "identity", "audit database", "audit", "server"}))
			Expect(fixture.handler).NotTo(BeNil())
			call := func(path string) int {
				w := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodGet, path, nil)
				if path != "/missing" {
					request.Pattern = "GET " + path
				}

				fixture.handler.ServeHTTP(w, request)
				return w.Code
			}
			Expect(call("/health")).To(Equal(http.StatusOK))
			Expect(call("/missing")).NotTo(Equal(http.StatusOK))
			fixture.draining.Store(true)
			Expect(call("/ready")).To(Equal(http.StatusServiceUnavailable))
			fixture.draining.Store(false)
			fixture.audit.EXPECT().Ready().Return(false)
			Expect(call("/ready")).To(Equal(http.StatusServiceUnavailable))
			fixture.audit.EXPECT().Ready().Return(true).Times(4)
			fixture.businessMock.ExpectPing().WillReturnError(errors.New("business unavailable"))
			Expect(call("/ready")).To(Equal(http.StatusServiceUnavailable))
			fixture.businessMock.ExpectPing()
			fixture.auditMock.ExpectPing().WillReturnError(errors.New("audit unavailable"))
			Expect(call("/ready")).To(Equal(http.StatusServiceUnavailable))
			fixture.businessMock.ExpectPing()
			fixture.auditMock.ExpectPing()
			fixture.messaging.EXPECT().Health(gomock.Any()).Return(errors.New("broker unavailable"))
			Expect(call("/ready")).To(Equal(http.StatusServiceUnavailable))
			fixture.businessMock.ExpectPing()
			fixture.auditMock.ExpectPing()
			fixture.messaging.EXPECT().Health(gomock.Any()).Return(nil)
			Expect(call("/ready")).To(Equal(http.StatusOK))
			Expect(fixture.businessMock.ExpectationsWereMet()).To(Succeed())
			Expect(fixture.auditMock.ExpectationsWereMet()).To(Succeed())
			fixture.server.EXPECT().Run(gomock.Any()).Return(nil)
			Expect(runtime.Run()).To(Succeed())
			runtime.Close()
			runtime.Close()
			Expect(fixture.closed).To(Equal([]string{"audit", "audit database", "identity", "integration", "foundation", "infrastructure"}))
		})
		It("rejects malformed required-tenant configuration before opening resources", func() {
			GinkgoT().Setenv("PLATFORM_TENANT_ACCESS_REQUIRED", "invalid")
			runtime, err := impl.New(nil).Build()
			Expect(runtime).To(BeNil())
			Expect(err).To(MatchError(ContainSubstring("load launcher configuration")))
		})
	})
}
