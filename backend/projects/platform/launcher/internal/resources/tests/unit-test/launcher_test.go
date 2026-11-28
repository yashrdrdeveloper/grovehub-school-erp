package launcher_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
	"grovehub/launcher"
)

// @Testing suite=LauncherFacade
func LauncherSpecs(){
	Describe("Public Launcher facade",func(){
		It("exposes the production logger, database, tenant and gateway composition boundary",func(){
			t:=GinkgoT()
			t.Setenv("GROVEHUB_LOG_LEVELS","")
			t.Setenv("DATABASE_URL","")
			t.Setenv("AUDIT_DATABASE_URL","")
			t.Setenv("PLATFORM_TENANT_ACCESS_REQUIRED","false")
			t.Setenv("PLATFORM_TENANT_ACCESS_ORIGIN","")
			t.Setenv("PLATFORM_TENANT_ACCESS_KEY_ID","")
			t.Setenv("PLATFORM_TENANT_ACCESS_PRIVATE_KEY","")
			logger,err:=launcher.NewLogger(&bytes.Buffer{})
			Expect(err).NotTo(HaveOccurred())
			Expect(logger).NotTo(BeNil())
			Expect(launcher.Run(logger)).To(MatchError(ContainSubstring("bootstrap PostgreSQL")))
			_,err=launcher.NewPostgresDB(context.Background())
			Expect(err).To(MatchError(launcher.ErrDatabaseURLRequired))
			shared:=&gorm.DB{Config:&gorm.Config{}}
			db,cleanup,err:=launcher.NewAuditPostgresDB(context.Background(),shared)
			Expect(err).NotTo(HaveOccurred())
			Expect(db).To(BeIdenticalTo(shared))
			cleanup()
			lookup,err:=launcher.NewTenantState(launcher.TenantAccessConfig{})
			Expect(err).NotTo(HaveOccurred())
			Expect(lookup).To(BeNil())
			verifier,err:=launcher.NewWorkloadVerifier("",nil,nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(verifier).NotTo(BeNil())
			mux:=http.NewServeMux()
			mux.HandleFunc("GET /health",func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(204)})
			gateway:=launcher.NewExternalGateway(mux,nil,verifier,"")
			response:=httptest.NewRecorder()
			gateway.ServeHTTP(response,httptest.NewRequest("GET","/health",nil))
			Expect(response.Code).To(Equal(204))
		})
	})
}
