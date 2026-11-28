package routes_test

import (
 "os"
 "path/filepath"
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "grovehub/libs/core/annotation/internal/src/generators/routes"
)

var _ = Describe("Launcher annotations",func(){
 It("generates only explicit edge probes without weakening service tiers",func(){
  t:=GinkgoT()
  backend,_:=testBackend(t)
  controller:=filepath.Join(backend,"launcher","internal","src","controller","health.go")
  writeTestFile(t,controller,`package controller
import "net/http"
// @Routes /
type HealthController interface {
 // @GET /health
 Health(http.ResponseWriter,*http.Request)
 // @GET /ready
 Readiness(http.ResponseWriter,*http.Request)
}
`)
  out:=filepath.Join(backend,"launcher",".generated","internal","src","routes","zz_routes.gen.go")
  generated:=generatedRoutes(t,backend,out)
  Expect(generated).To(ContainSubstring("grovehub/launcher/internal/src/controller"))
  Expect(generated).To(ContainSubstring("GET /health"))
  Expect(generated).To(ContainSubstring("GET /ready"))
  Expect(generated).To(ContainSubstring("registerHealthRoutes"))
  writeTestFile(t,controller,`package controller
import "net/http"
// @Routes /
type HealthController interface {
 // @POST /admin
 Admin(http.ResponseWriter,*http.Request)
}
`)
  Expect(routes.GenerateAll(backend)).To(MatchError(ContainSubstring("unsupported launcher route")))
  writeTestFile(t,controller,"not a Go file")
  Expect(routes.GenerateAll(backend)).NotTo(Succeed())
 })
 It("rejects an inaccessible Launcher controller root",func(){
  t:=GinkgoT()
  backend,_:=testBackend(t)
  dir:=filepath.Join(backend,"launcher","internal","src")
  Expect(os.MkdirAll(dir,0755)).To(Succeed())
  link:=filepath.Join(dir,"controller")
  Expect(os.Symlink("controller",link)).To(Succeed())
  Expect(routes.GenerateAll(backend)).NotTo(Succeed())
 })
})
