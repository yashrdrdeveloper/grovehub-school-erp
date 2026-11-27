package gateway

import (
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "strings"

 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 school "grovehub/libs/schoolgateway"
 integration "grovehub/services/integration/client/principal"
)

const gmailWebhook = "POST /integration/restricted/webhooks/gmail"

type delegatedIdentity interface {
 VerifyDelegatedSession(context.Context,string)(school.Principal,error)
 school.ProtectedRouteAuthorizer
 school.AccessVerifier
}

type platformWorkloads struct {
 scoped core.ScopedWorkloadVerifier
 policies map[string]core.RoutePolicy
 identity delegatedIdentity
 webhook core.WorkloadVerifier
}

type workloadConfig struct {
 Keys []core.WorkloadKeyConfig
 Routes map[string]core.RoutePolicy
}

// NewPlatformWorkloads uses explicit route grants. Empty config grants no
// service-to-service access. The existing Gmail webhook verifier stays isolated.
func NewPlatformWorkloads(raw string, identity delegatedIdentity, webhook core.WorkloadVerifier) (core.WorkloadVerifier,error) {
 v:=&platformWorkloads{identity:identity,webhook:webhook}
 if strings.TrimSpace(raw)=="" { return v,nil }
 var cfg workloadConfig
 decoder:=json.NewDecoder(strings.NewReader(raw))
 decoder.DisallowUnknownFields()
 if err:=decoder.Decode(&cfg);err!=nil {return nil,fmt.Errorf("invalid workload configuration")}
 var extra any
 if err:=decoder.Decode(&extra);err!=io.EOF {return nil,errors.New("invalid trailing workload configuration")}
 credentials,err:=core.NewEd25519Credentials(cfg.Keys)
 if err!=nil {return nil,err}
 if len(cfg.Routes)==0 {return nil,errors.New("no workload route policies")}
 for pattern,policy:=range cfg.Routes {
  if policy.Scope=="" || strings.TrimSpace(policy.Scope)!=policy.Scope || pattern==gmailWebhook {
   return nil,errors.New("invalid workload route policy")
  }
  if policy.RequireActor {
   if identity==nil || (!strings.Contains(pattern," /identity/restricted/") && !strings.Contains(pattern," /audit/restricted/")) {
    return nil,errors.New("delegated route has no supported identity authorization boundary")
   }
  } else {
   switch pattern {
   case "POST /identity/restricted/auth/login", "POST /identity/restricted/auth/forgot-password",
    "POST /identity/restricted/auth/refresh",
    "GET /identity/restricted/auth/recovery/validate",
    "GET /identity/restricted/auth/invitations/validate",
    "POST /identity/restricted/auth/recovery/exchange",
    "POST /identity/restricted/auth/invitations/exchange",
    "GET /identity/restricted/auth/reset-password/status",
    "POST /identity/restricted/auth/reset-password",
    "GET /identity/restricted/accounts/setup/status",
    "POST /identity/restricted/accounts/setup":
   default: return nil,errors.New("service-only route is not enabled for migration")
   }
  }
 }
 // The proof-based School boundary independently enforces actor presence.
 // The generic verifier must not require an unverified actor assertion.
 routes:=make(map[string]core.RoutePolicy,len(cfg.Routes))
 for pattern,policy:=range cfg.Routes {routes[pattern]=core.RoutePolicy{Scope:policy.Scope}}
 v.scoped=core.ScopedWorkloadVerifier{Credentials:credentials,Audience:"platform",Routes:routes}
 v.policies=cfg.Routes
 return v,nil
}

func (v *platformWorkloads) VerifyWorkload(ctx context.Context,r *http.Request)(context.Context,error) {
 if r==nil {return nil,core.ErrWorkloadDenied}
 if r.Pattern==gmailWebhook {
  if v.webhook==nil {return nil,core.ErrWorkloadDenied}
  return v.webhook.VerifyWorkload(ctx,r)
 }
 verified,err:=v.scoped.VerifyWorkload(ctx,r)
 if err!=nil {return nil,err}
 claims,ok:=core.WorkloadFrom(verified)
 if !ok {return nil,core.ErrWorkloadDenied}
 policy:=v.policies[r.Pattern]
 if !policy.RequireActor {
  if claims.Subject!="school" || claims.ActorGID!="" || claims.TenantID!="" || claims.SessionID!="" {return nil,core.ErrWorkloadDenied}
  return verified,nil
 }
 if v.identity==nil || claims.Subject!="school" {return nil,core.ErrWorkloadDenied}
 // School must prove possession of the actual access secret. A session ID
 // or a signed assertion alone cannot authenticate a browser user.
 values:=r.Header.Values("X-GH-Identity-Proof")
 if len(values)!=1 || len(values[0])==0 || len(values[0])>2048 || strings.ContainsAny(values[0]," \t\r\n,;") {return nil,core.ErrWorkloadDenied}
 if claims.SessionID!="" || claims.ActorGID!="" || claims.TenantID!="" {return nil,core.ErrWorkloadDenied}
 actor,err:=v.identity.VerifyAccess(verified,values[0])
 if err!=nil || actor.SessionID=="" || actor.AccountID=="" {return nil,core.ErrWorkloadDenied}
 r.Header.Del("X-GH-Identity-Proof")
 verified=school.WithPrincipal(verified,actor)
 // Reuse existing route policy during migration; unknown routes deny.
 oldPattern:=strings.Replace(r.Pattern,"/restricted/","/protected/",1)
 if err:=v.identity.AuthorizeProtectedRoute(verified,oldPattern);err!=nil {return nil,core.ErrWorkloadDenied}
 verified=integration.WithVerifiedPrincipal(verified,integration.VerifiedPrincipal{TenantID:actor.TenantID,ActorGID:actor.PersonGID})
 return verified,nil
}
