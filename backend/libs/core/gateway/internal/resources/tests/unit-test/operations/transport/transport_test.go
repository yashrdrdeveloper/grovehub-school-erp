package transport_test

import (
	"crypto/ed25519"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/operations/transport"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type roundTripFunc func(*http.Request) (*http.Response,error)
func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response,error) { return f(r) }

var _ = Describe("WorkloadTransport", func() {
	var value *transport.WorkloadTransport
	var request *http.Request
	BeforeEach(func() {
		_, key, err := ed25519.GenerateKey(nil)
		Expect(err).NotTo(HaveOccurred())
		origin, err := url.Parse("https://platform.example")
		Expect(err).NotTo(HaveOccurred())
		request, err = http.NewRequest("GET","https://platform.example/service/restricted/read",nil)
		Expect(err).NotTo(HaveOccurred())
		value = &transport.WorkloadTransport{Origin:origin, KeyID:"key", Issuer:"issuer", Subject:"subject", Audience:"platform", PrivateKey:key, Policy:func(*http.Request)(contracts.WorkloadClaims,error){return contracts.WorkloadClaims{Scopes:[]string{"read"}},nil}}
	})

	It("signs only allowed headers on a clone and disables redirects", func() {
		request.Header.Set("Cookie","secret")
		request.Header.Set("Authorization","browser")
		request.Header.Set("Idempotency-Key","command")
		value.Base = roundTripFunc(func(r *http.Request)(*http.Response,error){
			Expect(r.Header.Get("Cookie")).To(BeEmpty())
			Expect(r.Header.Get("Idempotency-Key")).To(Equal("command"))
			Expect(r.Header.Get("Authorization")).To(HavePrefix("Bearer ghw1."))
			return &http.Response{StatusCode:302,Header:http.Header{"Location":[]string{"https://other.example"}},Body:io.NopCloser(strings.NewReader("")),Request:r},nil
		})
		client, err := transport.NewWorkloadHTTPClient(value)
		Expect(err).NotTo(HaveOccurred())
		response, err := client.Do(request)
		Expect(err).NotTo(HaveOccurred())
		defer response.Body.Close()
		Expect(response.StatusCode).To(Equal(302))
		Expect(request.Header.Get("Authorization")).To(Equal("browser"))
	})

	It("rejects invalid dependencies, origins, host overrides and policies", func() {
		for _, bad := range []*transport.WorkloadTransport{nil, {}, {Origin:value.Origin}} {
			_, err := transport.NewWorkloadHTTPClient(bad)
			Expect(err).To(HaveOccurred())
			_, err = bad.RoundTrip(request)
			Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		}
		for _, address := range []string{"https://other.example", "http://platform.example", "https://user@platform.example"} {
			r, err := http.NewRequest("GET",address,nil)
			Expect(err).NotTo(HaveOccurred())
			_, err = value.RoundTrip(r)
			Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		}
		request.Host = "other.example"
		_, err := value.RoundTrip(request)
		Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
		request.Host = "platform.example"
		sentinel := errors.New("policy denied")
		value.Policy = func(*http.Request)(contracts.WorkloadClaims,error){return contracts.WorkloadClaims{},sentinel}
		_, err = value.RoundTrip(request)
		Expect(err).To(BeIdenticalTo(sentinel))
		value.Policy = func(*http.Request)(contracts.WorkloadClaims,error){return contracts.WorkloadClaims{},nil}
		_, err = value.RoundTrip(request)
		Expect(err).To(BeIdenticalTo(contracts.ErrWorkloadDenied))
	})
})
