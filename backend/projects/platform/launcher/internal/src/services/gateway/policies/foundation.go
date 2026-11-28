package policies

import core "grovehub/libs/core/gateway"

func foundationPolicies() map[string]core.RoutePolicy {
	return map[string]core.RoutePolicy{
		"POST /foundation/restricted/persons":                   {Scope: "foundation.persons.create", RequireActor: false},
		"GET /foundation/restricted/persons/{gid}":              {Scope: "foundation.persons.get", RequireActor: false},
		"PATCH /foundation/restricted/persons/{gid}":            {Scope: "foundation.persons.update", RequireActor: false},
		"POST /foundation/restricted/persons/search":            {Scope: "foundation.persons.search", RequireActor: false},
		"PATCH /foundation/restricted/persons/{gid}/deactivate": {Scope: "foundation.persons.deactivate", RequireActor: false},
	}
}
