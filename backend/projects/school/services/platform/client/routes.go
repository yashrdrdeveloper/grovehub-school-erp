// Package client defines the explicit School business-operation allowlist.
// Wire DTOs will be generated from verified service OpenAPI specifications.
package client

import (
 "net/http"
 "strings"
 auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/identity/client/authentication"
)

type Operation struct { Method, Path, Scope string }

var operations=[]Operation{
 {"POST","/identity/protected/accounts","identity.accounts.create"},
 {"POST","/identity/protected/accounts/search","identity.accounts.search"},
 {"GET","/identity/protected/accounts/roles","identity.accounts.roles"},
 {"GET","/identity/protected/accounts/{accountId}","identity.accounts.get"},
 {"PATCH","/identity/protected/accounts/{accountId}","identity.accounts.update"},
 {"PATCH","/identity/protected/accounts/{accountId}/suspend","identity.accounts.suspend"},
 {"PATCH","/identity/protected/accounts/{accountId}/activate","identity.accounts.activate"},
 {"PATCH","/identity/protected/accounts/{accountId}/deactivate","identity.accounts.deactivate"},
 {"POST","/identity/protected/accounts/{accountId}/resend-invitation","identity.accounts.resend"},
 {"GET","/identity/protected/account-roles/{accountId}","identity.accountroles.get"},
 {"POST","/identity/protected/account-roles/{accountId}/preview","identity.accountroles.preview"},
 {"POST","/identity/protected/approvals/history/search","identity.approvals.history"},
 {"POST","/identity/protected/approvals/inbox/search","identity.approvals.inbox"},
 {"GET","/identity/protected/approvals/inbox/count","identity.approvals.count"},
 {"GET","/identity/protected/approvals/{requestId}","identity.approvals.get"},
 {"POST","/identity/protected/approvals/lifecycle","identity.approvals.lifecycle"},
 {"POST","/identity/protected/approvals/permissions","identity.approvals.permissions"},
 {"POST","/identity/protected/approvals/assignments","identity.approvals.assignments"},
 {"POST","/identity/protected/approvals/{requestId}/approve","identity.approvals.approve"},
 {"POST","/identity/protected/approvals/{requestId}/reject","identity.approvals.reject"},
 {"POST","/identity/protected/approvals/{requestId}/cancel","identity.approvals.cancel"},
 {"POST","/identity/protected/approvals/{requestId}/execute","identity.approvals.execute"},
 {"GET","/identity/protected/capabilities/me","identity.capabilities.me"},
 {"GET","/identity/protected/capabilities/permissions","identity.capabilities.permissions"},
 {"POST","/identity/protected/roles","identity.roles.create"},
 {"PATCH","/identity/protected/roles/{roleId}","identity.roles.update"},
 {"POST","/identity/protected/roles/search","identity.roles.search"},
 {"GET","/identity/protected/roles/{roleId}","identity.roles.get"},
 {"GET","/identity/protected/roles/{roleId}/permissions","identity.roles.permissions"},
 {"POST","/identity/protected/roles/{roleId}/permissions/preview","identity.roles.preview"},
 {"POST","/audit/protected/search","audit.search"},
}

func Operations() []Operation { return append([]Operation(nil),operations...) }

func (o Operation) RestrictedPath() string { return strings.Replace(o.Path,"/protected/","/restricted/",1) }

// Resolve accepts only a registered method and path, with UUID-shaped IDs.
func (o Operation) Resolve(r *http.Request)(string,bool) {
 if r.Method!=o.Method{return "",false}
 parts:=strings.Split(o.Path,"/")
 actual:=strings.Split(r.URL.Path,"/")
 if len(parts)!=len(actual){return "",false}
 for i,part:=range parts{
  if strings.HasPrefix(part,"{") && strings.HasSuffix(part,"}"){
   if !auth.ValidSessionID(actual[i]){return "",false}
   continue
  }
  if part!=actual[i]{return "",false}
 }
 return strings.Replace(r.URL.Path,"/protected/","/restricted/",1),true
}

func Scope(r *http.Request)(string,bool){
 // Static routes take precedence over variable routes; UUID validation
 // prevents a static segment from being interpreted as an identifier.
 clone:=r.Clone(r.Context())
 u:=*r.URL
 u.Path=strings.Replace(u.Path,"/restricted/","/protected/",1)
 if u.Path==r.URL.Path{return "",false}
 clone.URL=&u
 for _,operation:=range operations{
  if _,ok:=operation.Resolve(clone);ok{return operation.Scope,true}
 }
 return "",false
}
