package launcher

import (
	"context"
	"errors"
	"fmt"
	auditclient "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/client/audit"
	auditcontroller "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/controller/audit"
	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/auditdelivery"
	corehttp "grovehub/libs/core/http"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"time"

	platform "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/client/platform"
	controller "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/controller/authentication"
	workflowController "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/controller/workflow"
	auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/authentication"
	workflow "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/resources/migrations"
	coredatabase "grovehub/libs/core/database"
)

// Run starts the private application listener with bounded readiness checks.
func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	db, err := coredatabase.OpenProjectPool(ctx, os.Getenv("INTERNAL_TOOL_DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("INTERNAL_TOOL_DATABASE_URL: %w", err)
	}
	defer db.Close()
	addr := os.Getenv("GROOVEHUB_HTTP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8082"
	}
	check, cancel := context.WithTimeout(ctx, 5*time.Second)
	err = migrations.Check(check, db)
	cancel()
	if err != nil {
		return err
	}
	authController, err := controller.New(auth.New(db), os.Getenv("INTERNAL_TOOL_ALLOWED_ORIGIN"))
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	authController.Register(mux)
	bindings, err := platform.Bindings(platform.Config{
		EnvironmentCode: os.Getenv("INTERNAL_TOOL_CATALOG_ENVIRONMENT"),
		TemplateID:      os.Getenv("INTERNAL_TOOL_CATALOG_TEMPLATE_ID"),
		Origin:          os.Getenv("INTERNAL_TOOL_PLATFORM_ORIGIN"),
		KeyID:           os.Getenv("INTERNAL_TOOL_WORKLOAD_KEY_ID"),
		PrivateKey:      os.Getenv("INTERNAL_TOOL_WORKLOAD_PRIVATE_KEY"),
	})
	if err != nil {
		return err
	}
	schoolConfig := platform.SchoolConfig{
		EnvironmentCode: os.Getenv("INTERNAL_TOOL_SCHOOL_ENVIRONMENT"),
		TemplateID:      os.Getenv("INTERNAL_TOOL_SCHOOL_TEMPLATE_ID"),
		SchoolID:        os.Getenv("INTERNAL_TOOL_SCHOOL_ID"),
		Origin:          os.Getenv("INTERNAL_TOOL_SCHOOL_ORIGIN"),
		KeyID:           os.Getenv("INTERNAL_TOOL_SCHOOL_KEY_ID"),
		PrivateKey:      os.Getenv("INTERNAL_TOOL_SCHOOL_PRIVATE_KEY"),
		OperatorIDs:     os.Getenv("INTERNAL_TOOL_SCHOOL_OPERATOR_IDS"),
	}
	directoryConfig := platform.DirectoryConfig{Origin: os.Getenv("INTERNAL_TOOL_DIRECTORY_ORIGIN"), KeyID: os.Getenv("INTERNAL_TOOL_DIRECTORY_KEY_ID"), PrivateKey: os.Getenv("INTERNAL_TOOL_DIRECTORY_PRIVATE_KEY"), ReaderID: os.Getenv("INTERNAL_TOOL_DIRECTORY_READER_ID")}
	directory, err := platform.NewDirectory(directoryConfig)
	if err != nil {
		return err
	}
	receiptConfig := auditclient.Config{Origin: os.Getenv("INTERNAL_TOOL_RECEIPT_READ_ORIGIN"), KeyID: os.Getenv("INTERNAL_TOOL_RECEIPT_READ_KEY_ID"), PrivateKey: os.Getenv("INTERNAL_TOOL_RECEIPT_READ_PRIVATE_KEY")}
	receiptReader, err := auditclient.NewReceiptReader(receiptConfig)
	if err != nil {
		return err
	}
	if receiptReader != nil && directory != nil && (receiptConfig.KeyID == directoryConfig.KeyID || receiptConfig.PrivateKey == directoryConfig.PrivateKey) {
		return errors.New("receipt reader key must be independent")
	}
	reserved := []platform.ReservedCredential{{KeyID: directoryConfig.KeyID, PrivateKey: directoryConfig.PrivateKey}, {KeyID: receiptConfig.KeyID, PrivateKey: receiptConfig.PrivateKey}}
	for _, prefix := range []string{"INTERNAL_TOOL_WORKLOAD", "INTERNAL_TOOL_AUDIT", "INTERNAL_TOOL_AUDIT_READ"} {
		if directory != nil && (directoryConfig.KeyID == os.Getenv(prefix+"_KEY_ID") || directoryConfig.PrivateKey == os.Getenv(prefix+"_PRIVATE_KEY")) {
			return errors.New("school directory key must be independent")
		}
		if receiptReader != nil && (receiptConfig.KeyID == os.Getenv(prefix+"_KEY_ID") || receiptConfig.PrivateKey == os.Getenv(prefix+"_PRIVATE_KEY")) {
			return errors.New("receipt reader key must be independent")
		}
		reserved = append(reserved, platform.ReservedCredential{KeyID: os.Getenv(prefix + "_KEY_ID"), PrivateKey: os.Getenv(prefix + "_PRIVATE_KEY")})
	}
	schoolBindings, schoolOrigins, err := platform.ParseSchoolGrants(os.Getenv("INTERNAL_TOOL_SCHOOL_GRANTS"), schoolConfig, reserved, directoryConfig.Origin)
	if err != nil {
		return err
	}
	workflowController.New(workflow.NewWithSchoolBindings(db, bindings, schoolBindings).WithDirectory(directory), authController).Register(mux)
	publisher, err := auditclient.New(auditclient.Config{Origin: os.Getenv("INTERNAL_TOOL_AUDIT_ORIGIN"), KeyID: os.Getenv("INTERNAL_TOOL_AUDIT_KEY_ID"), PrivateKey: os.Getenv("INTERNAL_TOOL_AUDIT_PRIVATE_KEY")})
	if err != nil {
		return err
	}
	reader, err := auditclient.NewReader(auditclient.Config{Origin: os.Getenv("INTERNAL_TOOL_AUDIT_READ_ORIGIN"), KeyID: os.Getenv("INTERNAL_TOOL_AUDIT_READ_KEY_ID"), PrivateKey: os.Getenv("INTERNAL_TOOL_AUDIT_READ_PRIVATE_KEY")})
	if err != nil {
		return err
	}
	if reader != nil && publisher != nil && (os.Getenv("INTERNAL_TOOL_AUDIT_READ_KEY_ID") == os.Getenv("INTERNAL_TOOL_AUDIT_KEY_ID") || os.Getenv("INTERNAL_TOOL_AUDIT_READ_PRIVATE_KEY") == os.Getenv("INTERNAL_TOOL_AUDIT_PRIVATE_KEY")) {
		return errors.New("audit reader and writer must use independent keys")
	}
	monitor := auditdelivery.NewMonitor(db, publisher != nil)
	auditcontroller.Register(mux, authController, monitor, reader, receiptReader)
	mux.HandleFunc("GET /health/audit-delivery", monitor.HealthHandler)
	if publisher != nil {
		workerCtx, workerStop := context.WithCancel(ctx)
		workerDone := make(chan struct{})
		go func() { defer close(workerDone); auditdelivery.New(db, publisher).Run(workerCtx) }()
		defer func() { workerStop(); <-workerDone }()
	}
	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	mux.HandleFunc("GET /health/database", func(w http.ResponseWriter, r *http.Request) {
		probe, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(probe); err != nil {
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	var draining atomic.Bool
	checks := []corehttp.Check{db.PingContext}
	checks = append(checks, func(ctx context.Context) error {
		if err := migrations.Check(ctx, db); err != nil {
			return err
		}
		var active bool
		if err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM internal_tool.operators WHERE persona='L2' AND status='ACTIVE')").Scan(&active); err != nil {
			return err
		}
		if !active {
			return errors.New("no active L2 operator")
		}
		return nil
	})
	if len(bindings) > 0 {
		upstream, err := corehttp.ReadinessHTTP(os.Getenv("INTERNAL_TOOL_PLATFORM_ORIGIN"), "/ready")
		if err != nil {
			return err
		}
		checks = append(checks, upstream)
	}
	for _, origin := range schoolOrigins {
		if len(bindings) > 0 && origin == os.Getenv("INTERNAL_TOOL_PLATFORM_ORIGIN") {
			continue
		}
		upstream, err := corehttp.ReadinessHTTP(origin, "/ready")
		if err != nil {
			return err
		}
		checks = append(checks, upstream)
	}
	mux.HandleFunc("GET /health/ready", corehttp.ReadinessHandler(&draining, checks...))
	server := &http.Server{Addr: addr, Handler: corehttp.Ingress(mux), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	done := make(chan error, 1)
	go func() { done <- server.ListenAndServe() }()
	select {
	case err := <-done:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		draining.Store(true)
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdown)
	}
}
