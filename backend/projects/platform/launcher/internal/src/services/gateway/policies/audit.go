package policies

import core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"

func auditPolicies() map[string]core.RoutePolicy {
	return map[string]core.RoutePolicy{
		"POST /audit/restricted/search": {Scope: "audit.search", RequireActor: true},
	}
}
