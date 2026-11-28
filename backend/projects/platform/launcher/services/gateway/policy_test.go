package gateway

import (
 "testing"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)

func TestRestrictedPolicyCannotBeWeakened(t *testing.T){
 cases:=[]struct{pattern,scope string;actor,allowed bool}{
 {"POST /identity/restricted/auth/login","identity.auth.login",false,true},
 {"POST /identity/restricted/auth/login","identity.auth.login",true,false},
 {"GET /identity/restricted/auth/me","identity.auth.me",true,true},
 {"GET /identity/restricted/auth/me","identity.auth.me",false,false},
 {"GET /identity/restricted/auth/me","identity.auth.login",true,false},
 {"GET /identity/restricted/accounts/{accountId}/person","identity.accounts.person",true,true},
 {"GET /identity/restricted/accounts/{accountId}/person","identity.accounts.person",false,false},
 {"GET /identity/restricted/accounts/{accountId}/person","identity.accounts.get",true,false},
 {"POST /audit/restricted/search","audit.search",true,true},
 {"POST /audit/restricted/search","audit.search",false,false},
 {"POST /integrations/restricted/webhooks/gmail","integration.webhook",false,false},
 {"POST /foundation/restricted/persons","foundation.persons.create",false,true},
 {"POST /foundation/restricted/persons","foundation.persons.create",true,false},
 {"GET /identity/restricted/unknown","identity.auth.me",true,false},
 }
 for _,tc:=range cases {if got:=validRoutePolicy(tc.pattern,core.RoutePolicy{Scope:tc.scope,RequireActor:tc.actor});got!=tc.allowed{t.Errorf("%s: got %v",tc.pattern,got)}}
}
