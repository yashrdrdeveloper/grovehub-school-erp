// Package readiness is a compatibility facade.
// Deprecated: use the root core/http package for new code.
package readiness

import ("net/http"; "sync/atomic"; corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http")
type Check = corehttp.Check
func Handler(draining *atomic.Bool, checks ...Check) http.HandlerFunc { return corehttp.ReadinessHandler(draining, checks...) }
func HTTP(origin,path string)(Check,error){return corehttp.ReadinessHTTP(origin,path)}
