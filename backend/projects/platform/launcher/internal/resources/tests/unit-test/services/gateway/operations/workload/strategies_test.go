package workload_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	fixture "grovehub/launcher/internal/resources/tests/helpers/workload"
	"grovehub/launcher/internal/src/services/gateway/operations/workload"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/internalcatalog"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/operatoraudit"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/operatorauditread"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/schoolexecution"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/schoolreceipt"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/schoolreceiptread"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/strategies/schoolstatusreceipt"
)

type strategyCase struct {
	name, pattern, path, issuer, subject, scope string
	newVerifier func([]core.WorkloadKeyConfig, []core.WorkloadKeyConfig) (core.WorkloadVerifier, error)
	configure func(*workload.WorkloadConfig, []core.WorkloadKeyConfig)
}

func strategyCases() []strategyCase {
	school:=uuid.NewString()
	return []strategyCase{
		{"catalog",internalcatalog.InternalCatalogRoute,"/integrations/restricted/providers","internal-tool","internal-tool",internalcatalog.InternalCatalogScope,internalcatalog.New,func(c *workload.WorkloadConfig,k []core.WorkloadKeyConfig){c.InternalToolKeys=k}},
		{"audit writer",operatoraudit.OperatorAuditRoute,"/audit/restricted/operator-activities","internal-tool-audit","internal-tool-audit",operatoraudit.OperatorAuditScope,operatoraudit.New,func(c *workload.WorkloadConfig,k []core.WorkloadKeyConfig){c.OperatorAuditKeys=k}},
		{"audit reader",operatorauditread.OperatorAuditReadRoute,"/audit/restricted/operator-activities/search","internal-tool-audit-reader","internal-tool-audit-reader",operatorauditread.OperatorAuditReadScope,operatorauditread.New,func(c *workload.WorkloadConfig,k []core.WorkloadKeyConfig){c.OperatorAuditReadKeys=k}},
		{"school execution",schoolexecution.SchoolAuditCheckRoute,"/audit/restricted/schools/"+school+"/activity-check",schoolexecution.SchoolExecutionIssuer,schoolexecution.SchoolSubjectPrefix+school,schoolexecution.SchoolAuditCheckScope,schoolexecution.New,func(c *workload.WorkloadConfig,k []core.WorkloadKeyConfig){c.SchoolExecutionKeys=k}},
		{"rename receipt",schoolreceipt.ReceiptRoute,"/audit/restricted/school-rename-receipts",schoolreceipt.ReceiptIssuer,schoolreceipt.ReceiptIssuer+":"+school,schoolreceipt.ReceiptScope,schoolreceipt.New,func(c *workload.WorkloadConfig,k []core.WorkloadKeyConfig){c.SchoolReceiptKeys=k}},
		{"receipt reader",schoolreceiptread.ReceiptReadRoute,"/audit/restricted/school-mutation-receipts/search","internal-tool-school-receipt-reader","internal-tool-school-receipt-reader",schoolreceiptread.ReceiptReadScope,schoolreceiptread.New,func(c *workload.WorkloadConfig,k []core.WorkloadKeyConfig){c.SchoolReceiptReadKeys=k}},
		{"status receipt",schoolstatusreceipt.StatusReceiptRoute,"/audit/restricted/school-status-receipts",schoolstatusreceipt.StatusReceiptIssuer,schoolstatusreceipt.StatusReceiptIssuer+":"+school,schoolstatusreceipt.StatusReceiptScope,schoolstatusreceipt.New,func(c *workload.WorkloadConfig,k []core.WorkloadKeyConfig){c.SchoolStatusReceiptKeys=k}},
	}
}

// @Testing suite=LauncherWorkloads
func StrategyFailureSpecs() {
	Describe("Dedicated workload strategies",func(){
		It("rejects malformed, shared and invalid cryptographic configuration for every strategy",func(){
			for _,tc:=range strategyCases(){
				By(tc.name)
				key,_,err:=fixture.Key("strategy",tc.issuer,tc.subject,[]string{tc.scope})
				Expect(err).NotTo(HaveOccurred())
				for _,mutate:=range []func(*core.WorkloadKeyConfig){
					func(k *core.WorkloadKeyConfig){k.Subject="foreign"},
					func(k *core.WorkloadKeyConfig){k.AllowDelegation=true},
					func(k *core.WorkloadKeyConfig){k.Scopes=[]string{"wrong"}},
					func(k *core.WorkloadKeyConfig){k.PublicKey="not-a-key"},
				}{
					bad:=key
					mutate(&bad)
					_,err:=tc.newVerifier([]core.WorkloadKeyConfig{bad},nil)
					Expect(err).To(HaveOccurred())
					var cfg workload.WorkloadConfig
					tc.configure(&cfg,[]core.WorkloadKeyConfig{bad})
					raw,err:=json.Marshal(cfg)
					Expect(err).NotTo(HaveOccurred())
					_,err=workload.NewPlatformWorkloads(string(raw),nil,nil)
					Expect(err).To(HaveOccurred())
				}
				_,err=tc.newVerifier([]core.WorkloadKeyConfig{key},[]core.WorkloadKeyConfig{key})
				Expect(err).To(HaveOccurred())
				_,err=tc.newVerifier(nil,nil)
				Expect(err).To(HaveOccurred())
			}
		})
		It("denies wrong route shape, absent credentials and duplicate scopes while accepting exactly one trusted scope",func(){
			for _,tc:=range strategyCases(){
				By(tc.name)
				key,private,err:=fixture.Key("strategy",tc.issuer,tc.subject,[]string{tc.scope})
				Expect(err).NotTo(HaveOccurred())
				verifier,err:=tc.newVerifier([]core.WorkloadKeyConfig{key},nil)
				Expect(err).NotTo(HaveOccurred())
				method,_,_:=strings.Cut(tc.pattern," ")
				for _,mutate:=range []func(*http.Request){
					func(r *http.Request){r.URL=nil},
					func(r *http.Request){r.URL.RawQuery="unexpected=1"},
					func(r *http.Request){r.Header.Set("Cookie","untrusted=x")},
					func(r *http.Request){r.Header.Set("X-GH-Identity-Proof","untrusted")},
					func(r *http.Request){r.Header.Del("Authorization")},
					func(r *http.Request){r.Method="DELETE"},
					func(r *http.Request){r.URL.Path="/wrong"},
				}{
					r,err:=fixture.Request(tc.pattern,tc.path,key,private,core.WorkloadClaims{Scopes:[]string{tc.scope}})
					Expect(err).NotTo(HaveOccurred())
					mutate(r)
					_,err=verifier.VerifyWorkload(context.Background(),r)
					Expect(err).To(HaveOccurred())
				}
				for _,scopes:=range [][]string{{tc.scope,tc.scope},{tc.scope}}{
					r,err:=fixture.Request(tc.pattern,tc.path,key,private,core.WorkloadClaims{Scopes:scopes})
					Expect(err).NotTo(HaveOccurred())
					Expect(r.Method).To(Equal(method))
					ctx,err:=verifier.VerifyWorkload(context.Background(),r)
					if len(scopes)>1 {
						Expect(err).To(MatchError(core.ErrWorkloadDenied))
					}else{
						Expect(err).NotTo(HaveOccurred())
						claims,ok:=core.WorkloadFrom(ctx)
						Expect(ok).To(BeTrue())
						Expect(claims.Subject).To(Equal(tc.subject))
					}
				}
			}
		})
		It("denies unconfigured special workload routes and cross-school credentials",func(){
			verifier,err:=workload.NewPlatformWorkloads("",nil,nil)
			Expect(err).NotTo(HaveOccurred())
			for _,tc:=range strategyCases(){
				method,_,_:=strings.Cut(tc.pattern," ")
				r:=httptest.NewRequest(method,tc.path,nil)
				r.Pattern=tc.pattern
				_,err:=verifier.VerifyWorkload(context.Background(),r)
				Expect(err).To(MatchError(core.ErrWorkloadDenied))
			}
			school:=uuid.NewString()
			key,private,err:=fixture.Key("school",schoolexecution.SchoolExecutionIssuer,schoolexecution.SchoolSubjectPrefix+school,[]string{schoolexecution.SchoolAuditCheckScope})
			Expect(err).NotTo(HaveOccurred())
			v,err:=schoolexecution.New([]core.WorkloadKeyConfig{key},nil)
			Expect(err).NotTo(HaveOccurred())
			r,err:=fixture.Request(schoolexecution.SchoolAuditCheckRoute,"/audit/restricted/schools/"+uuid.NewString()+"/activity-check",key,private,core.WorkloadClaims{Scopes:key.Scopes})
			Expect(err).NotTo(HaveOccurred())
			_,err=v.VerifyWorkload(context.Background(),r)
			Expect(err).To(MatchError(core.ErrWorkloadDenied))
		})
	})
}
