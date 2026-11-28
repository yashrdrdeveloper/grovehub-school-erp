package impl_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/launcher/internal/src/services/gateway/operations/ratelimit/impl"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

// @Testing suite=LauncherRateLimit
func RateLimitSpecs() {
	Describe("Allow", func() {
		It("bounds each configured route and caller independently", func() {
			limiter := impl.New(map[string]int{"POST /login": 2, "GET /disabled": 0})
			request := httptest.NewRequest(http.MethodPost, "/login", nil)
			request.RemoteAddr = "127.0.0.1:1000"
			Expect(limiter.Allow("/login", request)).To(BeTrue())
			request.RemoteAddr = "127.0.0.1:2000"
			Expect(limiter.Allow("/login", request)).To(BeTrue())
			Expect(limiter.Allow("/login", request)).To(BeFalse())
			request.RemoteAddr = "127.0.0.2:2000"
			Expect(limiter.Allow("/login", request)).To(BeTrue())
			Expect(limiter.Allow("/unknown", request)).To(BeTrue())
		})
		It("resets a caller after the window expires and accepts a bare remote address", func() {
			now := time.Unix(1700000000, 0)
			limiter := impl.NewWithClock(map[string]int{"POST /login": 1}, func() time.Time {
				return now
			})
			request := httptest.NewRequest(http.MethodPost, "/login", nil)
			request.RemoteAddr = "127.0.0.1"
			Expect(limiter.Allow("/login", request)).To(BeTrue())
			Expect(limiter.Allow("/login", request)).To(BeFalse())
			now = now.Add(time.Minute + time.Second)
			Expect(limiter.Allow("/login", request)).To(BeTrue())
		})
		It("fails closed when full, then removes expired entries before admitting a caller", func() {
			now := time.Unix(1700000000, 0)
			limiter := impl.NewWithClock(map[string]int{"POST /login": 1}, func() time.Time {
				return now
			})
			request := httptest.NewRequest(http.MethodPost, "/login", nil)
			for i := 0; i <= 100000; i++ {
				request.RemoteAddr = "caller-" + strconv.Itoa(i)
				if !limiter.Allow("/login", request) {
					Fail("limiter rejected a caller before its capacity guard")
				}
			}

			request.RemoteAddr = "overflow"
			Expect(limiter.Allow("/login", request)).To(BeFalse())
			now = now.Add(time.Minute + time.Second)
			Expect(limiter.Allow("/login", request)).To(BeTrue())
		})
	})
}
