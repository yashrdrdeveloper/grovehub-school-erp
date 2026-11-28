# Launcher Gateway — Engineering

## Purpose

The gateway is owned by `backend/launcher/services/gateway`, not Identity or Integration.

- `operations/external` classifies the **matched ServeMux route pattern** with an anchored regex; it never trusts the raw URL or client headers to determine the tier. After policy checks, it dispatches through `mux.ServeHTTP` so Go path variables are correctly populated.
- `public` permits unauthenticated requests, subject to rate limits, Origin checks, and CSRF on grant/refresh-cookie operations.
- `protected` validates Identity's opaque cookie session and attaches a trusted Principal to the request context. Integration receives a verified tenant/PersonGID through its public principal facade.
- `restricted` requires an explicit workload verifier, never a browser session. Gmail Pub/Sub push is verified against its configured audience and service account; all unknown workloads fail closed.
- `operations/inprocess` provides typed, in-process call authorization for verified caller context. The default authorizer is nil, which denies privileged internal dispatch until bootstrap explicitly supplies permissions.
- `/health` is the only path exempt from tier classification. Unknown/unclassified application routes fail closed.
- Services register annotated controller operations using generated route-group functions; no service implements a second HTTP auth middleware.

A process-local gateway throttle is not a distributed rate limit. Cluster deployments still require a shared edge limiter, a well-defined trusted-proxy configuration, and auth bootstrap. TLS termination and production Origin configuration must be explicit.


## RouteLimit annotation contract

The controller method declares both its route and its own requests-per-minute policy:

```go
// @POST /login
// @RouteLimit 12
Login(http.ResponseWriter, *http.Request)
```

`libs/core/annotation/internal/src/generators/routes` validates that a rate limit is an integer between 1 and 100000 and generates `routes.RouteLimits()` keyed by the exact method and registered URL pattern. Identity bootstrap exposes this policy through `Module.RouteLimits()` to Launcher. The generic limiter has no Identity path switch. Missing `@RouteLimit` means the route has no route-specific rate limit.

The current limiter is per process and per originating socket IP. In multi-replica deployments, replace this storage with a distributed limiter at the gateway boundary; only trust forwarding headers when a reverse proxy is explicitly trusted. This does not replace account-aware login throttling, lockout strategy, or operational abuse monitoring.

## Interface-first resource convention

Services expose public `Service` interfaces; `operations/<usecase>/<usecase>.go` exposes the use-case contract; `operations/<usecase>/impl` implements it; `services/<resource>/impl` composes those operations. Security verification and cookies have separate contracts and implementation packages. Generated routes are composed in `routes/http.go`.

The internal gateway is an authorization bridge for typed calls and does not inherit browser gateway authorization automatically. Missing internal authorizer denies privileged calls.
