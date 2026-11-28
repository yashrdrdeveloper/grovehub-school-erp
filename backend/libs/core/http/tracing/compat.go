// Package tracing is a compatibility facade.
// Deprecated: use the root core/http package for new code.
package tracing

import ("net/http"; "time"; corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http")
func NewTransport(base http.RoundTripper) http.RoundTripper {return corehttp.NewTransport(base)}
func NewServiceTransport(service string,base http.RoundTripper)http.RoundTripper{return corehttp.NewServiceTransport(service,base)}
func NewClient(timeout time.Duration)*http.Client{return corehttp.NewClient(timeout)}
func NewServiceClient(service string,timeout time.Duration)*http.Client{return corehttp.NewServiceClient(service,timeout)}
func ServiceHandler(service,operation string,next http.Handler)http.Handler{return corehttp.ServiceHandler(service,operation,next)}
