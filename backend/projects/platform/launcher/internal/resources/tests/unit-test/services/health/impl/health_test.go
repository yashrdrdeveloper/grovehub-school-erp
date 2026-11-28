package impl_test
import (
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "grovehub/launcher/internal/src/services/health/impl"
)
// @Testing suite=LauncherHealthService
func HealthServiceSpecs(){
 Describe("Readiness",func(){
  It("keeps liveness independent of dependency readiness",func(){
   service:=impl.New(nil)
   Expect(service.Health().Status).To(Equal("UP"))
   response,ready:=service.Readiness()
   Expect(response.Status).To(Equal("NOT_READY"))
   Expect(ready).To(BeFalse())
   service=impl.New(func()bool{return true})
   response,ready=service.Readiness()
   Expect(response.Status).To(Equal("READY"))
   Expect(ready).To(BeTrue())
  })
 })
}
