package policies

import core "grovehub/libs/core/gateway"

func auditPolicies() map[string]core.RoutePolicy {
	return map[string]core.RoutePolicy{
		"POST /audit/restricted/search": {Scope: "audit.search", RequireActor: true},
	}
}
