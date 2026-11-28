package helpers_test

import (
 "net/http"
 "net/http/httptest"
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "grovehub/launcher/internal/src/middleware/helpers"
)

// @Testing suite=LauncherResponseRecorder
func RecorderSpecs(){
 Describe("Recorder",func(){
  It("captures implicit status and body size, and preserves the first status",func(){
   base:=httptest.NewRecorder()
   recorder:=&helpers.Recorder{ResponseWriter:base}
   n,err:=recorder.Write([]byte("abc"))
   Expect(err).NotTo(HaveOccurred())
   Expect(n).To(Equal(3))
   recorder.WriteHeader(http.StatusForbidden)
   Expect(recorder.Status).To(Equal(http.StatusOK))
   Expect(recorder.Bytes).To(Equal(int64(3)))
   Expect(recorder.Unwrap()).To(BeIdenticalTo(base))
  })
 })
}
