package runtime_test

import (
	"context"
	"time"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/operations/runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metric recording without an active runtime",func(){
	It("is safe before initialization and after shutdown",func(){
		ctx:=context.Background()
		Expect(func(){
			runtime.RecordOperation(ctx,"identity","health",time.Millisecond,200)
			runtime.RecordHTTPServerRequest(ctx,"identity","GET","/health",time.Millisecond,500)
			runtime.AddActiveHTTPServerRequest(ctx,"identity","GET",1)
		}).NotTo(Panic())
	})
})
