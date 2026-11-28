// Package request is a compatibility facade.
// Deprecated: use the root core/http package for new code.
package request

import ("net/http"; corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http")
func DecodeJSON(req *http.Request,destination any)error{return corehttp.DecodeJSON(req,destination)}
func DecodeStrictJSON(req *http.Request,destination any)error{return corehttp.DecodeStrictJSON(req,destination)}
