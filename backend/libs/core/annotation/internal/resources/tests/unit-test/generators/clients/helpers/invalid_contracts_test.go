package helpers_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/clients/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Client discovery failures",func(){
	It("rejects invalid receivers, embedded contracts, conflicting directives and proxy bindings",func(){
		for _,test:=range []struct{source,reason string}{
			{"// @Client invalid-name\ntype Client interface{}","invalid @Client receiver"},
			{"// @Client impl\ntype Client interface{ Other }","cannot embed interfaces"},
			{"// @Client impl\ntype Client interface{ Read() }","declare exactly one"},
			{"// @Client impl\ntype Client interface{\n// @Proxy svc.Read\n// @RequestLine GET /records\nRead() }","declare exactly one"},
			{"// @Client impl\ntype Client interface{\n// @Proxy invalid\nRead() }","invalid @Proxy binding"},
			{"// @Client impl\n// @Client other\ntype Client interface{}","duplicate"},
			{"// @Client impl\ntype Client interface{\n// @Proxy svc.Read\n// @Proxy svc.Other\nRead() }","duplicate"},
		}{
			root:=GinkgoT().TempDir()
			fixture.Write(root,"client.go","package client\n"+test.source)
			_,err:=helpers.Discover(root,"")
			Expect(err).To(MatchError(ContainSubstring(test.reason)),test.source)
		}
	})
	It("rejects malformed request lines and incompatible method signatures",func(){
		for _,test:=range []struct{line,signature,reason string}{
			{"GET","(ctx context.Context, req Request) (Result,error)","requires METHOD"},
			{"GET records","(ctx context.Context, req Request) (Result,error)","invalid @RequestLine path"},
			{"GET /a/../records","(ctx context.Context, req Request) (Result,error)","unsafe @RequestLine path"},
			{"GET /records\n// @RequestLine GET /records","(ctx context.Context, req Request) (Result,error)","duplicate @RequestLine"},
			{"GET /records","(ctx context.Context) (Result,error)","requires context and one"},
			{"GET /records","(ctx string, req Request) (Result,error)","first @RequestLine parameter"},
			{"GET /records","(ctx context.Context, req []Request) (Result,error)","named request struct"},
			{"GET /records","(ctx context.Context, req Request) error","expected (response,error)"},
			{"GET /records","(ctx context.Context, req Request) (Result,string)","last result must be error"},
			{"GET /records","(ctx context.Context, req Request) (Result,string,error)","second result must be bool"},
			{"GET /records","(ctx context.Context, req Request) (Result,bool,error)","only for POST"},
			{"GET /records","(ctx context.Context, req Request) (out Result, err error)","named @RequestLine results"},
		}{
			root:=GinkgoT().TempDir()
			fixture.Write(root,"client.go","package client\n// @Client impl\ntype Client interface {\n// @RequestLine "+test.line+"\nRead"+test.signature+"\n}")
			routes:=fixture.Write(root,"routes/zz_routes.gen.go","package routes\nfunc register() { mux.Handle(\"GET /records\",nil) }")
			_,err:=helpers.Discover(root,routes)
			Expect(err).To(MatchError(ContainSubstring(test.reason)),test.signature)
		}
	})
	It("resolves exact routes while ignoring unrelated calls and unmarked declarations",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"ignored_test.go","not Go source")
		fixture.Write(root,"client.go","package client\ntype Plain interface{}\ntype Record struct{}\n// @Client impl\ntype Client interface {\n// @RequestLine GET /records\nRead(ctx context.Context, req Request) (Result,error)\n}")
		routes:=fixture.Write(root,"routes/zz_routes.gen.go","package routes\nfunc register() { helper(); mux.Other(); mux.Handle(variable,nil); mux.Handle(1,nil); mux.Handle(\"POST /records\",nil); mux.Handle(\"GET /records\",nil) }")
		proxies,err:=helpers.Discover(root,routes)
		Expect(err).NotTo(HaveOccurred())
		Expect(proxies).To(HaveLen(1))
		_,err=helpers.Discover(root,filepath.Join(root,"missing.go"))
		Expect(err).To(MatchError(ContainSubstring("read generated controller routes")))
	})
})
