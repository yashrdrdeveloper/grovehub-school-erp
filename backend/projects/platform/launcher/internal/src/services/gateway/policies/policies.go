package policies

import core "grovehub/libs/core/gateway"

var registry = build()

func build() map[string]core.RoutePolicy {
	result := make(map[string]core.RoutePolicy)
	for _, group := range []map[string]core.RoutePolicy{identityPolicies(), integrationPolicies(), foundationPolicies(), auditPolicies()} {
		for pattern, policy := range group {
			result[pattern] = policy
		}
	}

	return result
}

func Find(pattern string) (core.RoutePolicy, bool) {
	policy, ok := registry[pattern]
	return policy, ok
}

func Validate(pattern string, policy core.RoutePolicy) bool {
	expected, ok := Find(pattern)
	return ok && expected == policy
}
