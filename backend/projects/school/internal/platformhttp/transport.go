// Package platformhttp owns School's consumer-side HTTP transport mechanics.
package platformhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const MaxResponseBytes = 1 << 20

type Result struct {
	Status int
	Body []byte
	Cookies []*http.Cookie
}

type proofKey struct{}
type bodyKey struct{}

func WithIdentityProof(ctx context.Context, proof string) context.Context {
	return context.WithValue(ctx, proofKey{}, proof)
}

func WithBody(ctx context.Context, raw json.RawMessage) context.Context {
	return context.WithValue(ctx, bodyKey{}, append(json.RawMessage(nil),raw...))
}

// ProofTransport is used by the signed workload transport, not browser ingress.
type ProofTransport struct { Base http.RoundTripper }

func (t ProofTransport) RoundTrip(req *http.Request) (*http.Response,error) {
	clone := req.Clone(req.Context())
	clone.Header.Del("X-GH-Identity-Proof")
	clone.Header.Del("Cookie")
	if proof,ok:=req.Context().Value(proofKey{}).(string);ok && proof!="" {
		clone.Header.Set("X-GH-Identity-Proof",proof)
	}
	base:=t.Base
	if base==nil {base=http.DefaultTransport}
	return base.RoundTrip(clone)
}

// BoundedClient wraps the already configured workload client. Preserve raw
// business JSON before signing to avoid nullable or numeric normalization.
func BoundedClient(client *http.Client) *http.Client {
	bounded:=*client
	bounded.Transport=boundedTransport{base:client.Transport}
	return &bounded
}

type boundedTransport struct { base http.RoundTripper }
func(t boundedTransport) RoundTrip(req *http.Request)(*http.Response,error){
	if raw,ok:=req.Context().Value(bodyKey{}).(json.RawMessage);ok {
		clone:=req.Clone(req.Context())
		if req.Body!=nil { _=req.Body.Close() }
		clone.Body=io.NopCloser(bytes.NewReader(raw))
		clone.ContentLength=int64(len(raw))
		clone.GetBody=func()(io.ReadCloser,error){return io.NopCloser(bytes.NewReader(raw)),nil}
		clone.Header.Del("Content-Length")
		req=clone
	}
	base:=t.base
	if base==nil {base=http.DefaultTransport}
	resp,err:=base.RoundTrip(req)
	if err!=nil{return nil,err}
	if resp==nil||resp.Body==nil{return nil,fmt.Errorf("empty platform response")}
	defer resp.Body.Close()
	if resp.StatusCode<200||resp.StatusCode>=600||(resp.StatusCode>=300&&resp.StatusCode<400){return nil,fmt.Errorf("unexpected platform response")}
	body,err:=io.ReadAll(io.LimitReader(resp.Body,MaxResponseBytes+1))
	if err!=nil||len(body)>MaxResponseBytes||!json.Valid(body){return nil,fmt.Errorf("invalid platform response")}
	resp.Body=io.NopCloser(bytes.NewReader(body))
	return resp,nil
}

func SDKResult(resp *http.Response, callErr error)(Result,error){
	if resp==nil||resp.Body==nil{return Result{},fmt.Errorf("platform unavailable")}
	defer resp.Body.Close()
	if resp.StatusCode<300&&callErr!=nil{return Result{},fmt.Errorf("invalid platform success response")}
	body,err:=io.ReadAll(io.LimitReader(resp.Body,MaxResponseBytes+1))
	if err!=nil||len(body)>MaxResponseBytes||!json.Valid(body){return Result{},fmt.Errorf("invalid platform response")}
	return Result{Status:resp.StatusCode,Body:body,Cookies:resp.Cookies()},nil
}

func DecodeBody[T any](input json.RawMessage)(T,bool){
	var value T
	raw:=bytes.TrimSpace(input)
	if len(raw)==0||raw[0]!='{'||!json.Valid(raw)||json.Unmarshal(raw,&value)!=nil{return value,false}
	return value,true
}
