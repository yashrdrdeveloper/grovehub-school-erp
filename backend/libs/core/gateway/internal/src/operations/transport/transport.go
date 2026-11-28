package transport

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/constants"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/operations/credential"

 "crypto/ed25519"
 "fmt"
 "net/http"
 "net/url"
 "strings"
 "time"
)

// WorkloadTransport signs requests to one explicitly configured Platform origin.
// Callers must derive delegation from authenticated application context.
type WorkloadTransport struct {
 Origin *url.URL
 KeyID string
 Issuer string
 Subject string
 Audience string
 PrivateKey ed25519.PrivateKey
 Policy func(*http.Request) (contracts.WorkloadClaims, error)
 Base http.RoundTripper
}

func (t *WorkloadTransport) RoundTrip(req *http.Request) (*http.Response, error) {
 if t == nil || req == nil || req.URL == nil || t.Origin == nil || t.Policy == nil {
  return nil, contracts.ErrWorkloadDenied
 }
 if req.URL.Scheme != t.Origin.Scheme || !strings.EqualFold(req.URL.Host,t.Origin.Host) || req.URL.User != nil || (req.Host != "" && !strings.EqualFold(req.Host,t.Origin.Host)) {
  return nil, contracts.ErrWorkloadDenied
 }
 if req.URL.Scheme != "https" && !(req.URL.Scheme == "http" && (req.URL.Hostname()=="127.0.0.1" || req.URL.Hostname()=="::1")) {
  return nil, contracts.ErrWorkloadDenied
 }
 claims, err := t.Policy(req)
 if err != nil { return nil, err }
 if len(claims.Scopes)==0 { return nil, contracts.ErrWorkloadDenied }
 claims.Subject=t.Subject
 claims.Audience=t.Audience
 claims.ExpiresAt=time.Now().Add(constants.OutboundCredentialLifetime)
 token,err:=credential.SignWorkload(t.KeyID,t.Issuer,t.PrivateKey,claims)
 if err!=nil { return nil,err }
 clone:=req.Clone(req.Context())
 clone.Header=make(http.Header)
 for _, name:=range []string{"Accept","Content-Type","Idempotency-Key","Traceparent","Tracestate","X-Correlation-ID"} {
  for _, value:=range req.Header.Values(name) { clone.Header.Add(name,value) }
 }
 clone.Header.Set("Authorization","Bearer "+token)
 base:=t.Base
 if base==nil { base=http.DefaultTransport }
 return base.RoundTrip(clone)
}

// NewWorkloadHTTPClient disables redirects to prevent credential forwarding.
func NewWorkloadHTTPClient(t *WorkloadTransport) (*http.Client,error) {
 if t==nil || t.Origin==nil || t.Policy==nil || len(t.PrivateKey)!=ed25519.PrivateKeySize || t.KeyID=="" || t.Issuer=="" || t.Subject=="" || t.Audience=="" {
  return nil,fmt.Errorf("invalid workload transport configuration")
 }
 return &http.Client{Transport:t,Timeout:constants.HTTPClientTimeout,CheckRedirect:func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}},nil
}
