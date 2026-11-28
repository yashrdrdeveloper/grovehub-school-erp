package policies

import (
	core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	mail "grovehub/services/integration/client/mail"
)

func integrationPolicies() map[string]core.RoutePolicy {
	return map[string]core.RoutePolicy{
		"POST " + mail.ReconciliationPath:                                    {Scope: mail.ReconciliationScope, RequireActor: false},
		"POST /integrations/restricted/connections":                          {Scope: "integration.create", RequireActor: true},
		"GET /integrations/restricted/connections/{id}":                      {Scope: "integration.get", RequireActor: true},
		"GET /integrations/restricted/connections":                           {Scope: "integration.list", RequireActor: true},
		"DELETE /integrations/restricted/connections/{id}":                   {Scope: "integration.disconnect", RequireActor: true},
		"POST /integrations/restricted/connections/{id}/oauth/start":         {Scope: "integration.oauth.start", RequireActor: true},
		"POST /integrations/restricted/connections/{id}/test":                {Scope: "integration.test", RequireActor: true},
		"GET /integrations/restricted/connections/{id}/tests":                {Scope: "integration.tests", RequireActor: true},
		"POST /integrations/restricted/connections/{id}/watch":               {Scope: "integration.watch", RequireActor: true},
		"POST /integrations/restricted/connections/system-mail/sender":       {Scope: "integration.systemmail", RequireActor: true},
		"POST /integrations/restricted/connections/{id}/actions/{action}":    {Scope: "integration.action", RequireActor: true},
		"GET /integrations/restricted/providers":                             {Scope: "integration.providers", RequireActor: true},
		"GET /integrations/restricted/providers/{provider}":                  {Scope: "integration.provider", RequireActor: true},
		"GET /integrations/restricted/connections/oauth/{provider}/callback": {Scope: "integration.oauth.callback", RequireActor: false},
	}
}
