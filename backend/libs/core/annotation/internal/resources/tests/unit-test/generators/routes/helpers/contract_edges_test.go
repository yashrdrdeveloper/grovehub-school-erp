package helpers_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/routes/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Controller contract boundaries",func(){
	It("rejects invalid controller declarations, signatures and route metadata",func(){
		for _,test:=range []struct{source,reason string}{
			{"// @Routes records\ntype ItemController interface{}","invalid route path"},
			{"// @Routes /records\n// @Routes /other\ntype ItemController interface{}","duplicate @Routes"},
			{"// @Routes /records\ntype Item struct{}","requires a Controller interface"},
			{"// @Routes /records\ntype ItemController interface{ Other }","cannot embed interfaces"},
			{"// @Routes /records\ntype ItemController interface{ Read() }","must accept http.ResponseWriter"},
			{"// @Routes /records\ntype ItemController interface{ Read(http.ResponseWriter,http.Request) }","must accept http.ResponseWriter"},
			{"// @Routes /records\ntype ItemController interface{ Read(other.ResponseWriter,*http.Request) }","must accept http.ResponseWriter"},
			{"// @Routes /records\ntype ItemController interface{ Read(http.ResponseWriter,*http.Other) }","must accept http.ResponseWriter"},
			{"// @Routes /records\ntype ItemController interface{ Read(http.ResponseWriter,*other.Request) }","must accept http.ResponseWriter"},
			{"// @Routes /records\ntype ItemController interface{ Read(http.ResponseWriter,*http.Request) }","missing @GET"},
			{"// @Routes /records\ntype ItemController interface{\n// @GET /\n// @POST /\nRead(http.ResponseWriter,*http.Request) }","multiple HTTP verb"},
			{"// @Routes /records\ntype ItemController interface{\n// @GET\nRead(http.ResponseWriter,*http.Request) }","requires exactly one path"},
			{"// @Routes /records\ntype ItemController interface{\n// @GET /../bad\nRead(http.ResponseWriter,*http.Request) }","invalid route path"},
			{"// @Routes /records\ntype ItemController interface{\n// @GET /\n// @RouteLimit 1\n// @RouteLimit 2\nRead(http.ResponseWriter,*http.Request) }","duplicate @RouteLimit"},
			{"// @Routes /records\ntype ItemController interface{\n// @GET /\n// @RouteLimit\nRead(http.ResponseWriter,*http.Request) }","requires positive requests"},
			{"// @Routes /records\ntype ItemController interface{\n// @GET /\n// @RouteLimit 100001\nRead(http.ResponseWriter,*http.Request) }","must be 1..100000"},
		}{
			root:=GinkgoT().TempDir()
			fixture.Write(root,"protected/item.go","package protected\n"+test.source)
			_,err:=helpers.DiscoverWithTier(root,"/demo",true,false)
			Expect(err).To(MatchError(ContainSubstring(test.reason)),test.source)
		}
	})
	It("requires service prefixes for categorized controllers and handles invalid package paths",func(){
		root:=GinkgoT().TempDir()
		_,err:=helpers.DiscoverWithTier(filepath.Join(root,"missing"),"",false,false)
		Expect(err).To(HaveOccurred())
		source:="// @Routes /records\ntype ItemController interface{\n// @GET /\nRead(http.ResponseWriter,*http.Request) }"
		path:=fixture.Write(root,"protected/item.go","package wrong\n"+source)
		_,err=helpers.DiscoverWithTier(root,"/demo",true,false)
		Expect(err).To(MatchError(ContainSubstring("expected controller package")))
		fixture.Write(root,"protected/item.go","package protected\n"+source)
		_,err=helpers.DiscoverWithTier(root,"",false,false)
		Expect(err).To(MatchError(ContainSubstring("categorized controller requires")))
		Expect(os.Remove(path)).To(Succeed())
		fixture.Write(root,"item.go","package controller\n"+source)
		_,err=helpers.DiscoverWithTier(root,"/demo",true,false)
		Expect(err).To(MatchError(ContainSubstring("must live in public/")))
		_,err=helpers.DiscoverWithTier(root,"",false,true)
		Expect(err).To(MatchError(ContainSubstring("restricted project requires")))
	})
	It("validates service prefix declarations independently of controller discovery",func(){
		root:=GinkgoT().TempDir()
		path:=filepath.Join(root,"http.go")
		prefix,found,err:=helpers.ServicePrefix(path)
		Expect(err).NotTo(HaveOccurred());Expect(found).To(BeFalse());Expect(prefix).To(BeEmpty())
		for _,test:=range []struct{source,reason string}{
			{"// @Routes /demo\nfunc Other(){}","must annotate the RegisterHTTP"},
			{"// @Routes /{service}\nfunc RegisterHTTP(){}","cannot contain path variables"},
			{"// @Routes /demo\nfunc RegisterHTTP(){}\n// @Routes /other\nfunc RegisterHTTP(){}","duplicate service @Routes"},
			{"// @Routes demo\nfunc RegisterHTTP(){}","invalid route path"},
			{"func broken(","expected"},
		}{
			fixture.Write(root,"http.go","package routes\n"+test.source)
			_,_,err:=helpers.ServicePrefix(path)
			Expect(err).To(MatchError(ContainSubstring(test.reason)))
		}
	})
})
