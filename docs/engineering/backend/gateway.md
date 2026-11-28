# Gateway: External and Internal Bridges

**Status:** Architecture standard.

## External ingress

The launcher composes a single HTTP gateway. The gateway MUST classify a **registered**, HTTP method-specific `ServeMux` pattern using an anchored access-tier matcher. The access tiers are public, protected, and restricted. A route outside these tiers MUST NOT be exposed unless it is explicitly allowlisted as a health endpoint.

- Public: no access session required; validate request size, Origin for writes, appropriate CSRF, and rate limits.
- Protected: validate an authenticated, nonrevoked principal, enforce Origin/CSRF, inject trusted principal into context. Service operations still enforce authorization and tenancy.
- Restricted: validate explicit workload/provider identity. Never accept a browser cookie as workload credentials.
- In all cases, dispatch via `ServeMux.ServeHTTP` after the gateway check so Go path variables remain populated.

The gateway MUST NOT trust X-Forwarded-For without a configured trusted proxy. In-memory throttles are defense in depth, not distributed abuse prevention.

## Internal invocation

Inter-module calls use typed public client facades. `@Proxy` dispatch is in-process and DOES NOT automatically run the external HTTP gateway. The internal gateway operation must verify the caller context and grant the requested capability before invoking the typed operation. Unknown internal principals/permissions fail closed. Do not expose a generic capability that accepts an arbitrary caller identity from input JSON or headers.

## Deployment contract

Keep generated @GET/@POST/@PATCH controller annotations as the single source of truth for registered methods and resource paths. Define the service prefix once in `routes/http.go`; it should only invoke generated registrations. Generate routes and format changed Go files as part of CI; test exact path, tier, account ownership, Origin, CSRF, wildcard binding, workload verification, and structured errors.

A gateway is a *transport policy boundary*, not a replacement for application resource-level authorization.

## Contract ownership

School session identity, AccessVerifier, legacy Admin/PersonGID claims and browser cookie names belong to `backend/projects/platform/launcher/services/gateway/session`. Generic route tier and workload verifier contracts remain in `libs/backend/core/gateway`. Verified contexts are established only by the trusted school adapter. This separation does not alter school session or origin/CSRF behavior and does not implement school RBAC enforcement by itself. The internal platform must own a separate session contract and cookie namespace.
