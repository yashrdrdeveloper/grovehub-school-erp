package helpers_test

import (
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "grovehub/launcher/internal/src/bootstrap/helpers"
)
// @Testing suite=LauncherCleanup
func CleanupSpecs(){
 Describe("Cleanup",func(){
  It("releases resources in reverse order exactly once",func(){
   var values []int
   cleanup:=&helpers.Cleanup{}
   cleanup.Add(func(){values=append(values,1)})
   cleanup.Add(nil)
   cleanup.Add(func(){values=append(values,2)})
   cleanup.Close()
   cleanup.Close()
   Expect(values).To(Equal([]int{2,1}))
  })
 })
}
