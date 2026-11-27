# Launcher Gateway — Request Behavior

The public edge accepts only routes that have a registered method, a matching route pattern, and a supported access tier.

- **Public:** an anonymous request is eligible for routing, subject to Origin, CSRF when cookies authorize an operation, and any annotated rate limit.
- **Protected:** a valid Identity access session is required. Unsafe requests also require CSRF and permitted Origin. A controller must independently authorize resource ownership and tenant access.
- **Restricted:** a verified provider/workload identity is required; browser sessions do not qualify. Unknown workloads are denied.
- **Unclassified:** denied by default, except registered health endpoints.

A rate-limited client receives HTTP 429 and the shared structured error envelope with a `Retry-After` header. Exact thresholds are declared at controller methods with `@RouteLimit` and generated into the gateway policy.

Test behavior: public login with a valid Origin and password succeeds if the account is eligible; an expired session on a protected route returns 401; an unknown restricted workload returns 403; a request exceeding the route-specific per-minute threshold returns 429. These are expected contracts, not proof of live verification.
