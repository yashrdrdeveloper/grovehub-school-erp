# Gateway boundary

The gateway is owned by `backend/launcher/services/gateway`, not Identity or Integration.

- `operations/external` classifies the **matched ServeMux route pattern** with an anchored regex; it never trusts the raw URL or client headers to determine the tier. After policy checks, it dispatches through `mux.ServeHTTP` so Go path variables are correctly populated.
- `public` permits unauthenticated requests, subject to rate limits, Origin checks, and CSRF on grant/refresh-cookie operations.
- `protected` validates Identity's opaque cookie session and attaches a trusted Principal to the request context. Integration receives a verified tenant/PersonGID through its public principal facade.
- `restricted` requires an explicit workload verifier, never a browser session. Gmail Pub/Sub push is verified against its configured audience and service account; all unknown workloads fail closed.
- `operations/internal` provides typed, in-process call authorization for verified caller context. The default authorizer is nil, which denies privileged internal dispatch until bootstrap explicitly supplies permissions.
- `/health` is the only path exempt from tier classification. Unknown/unclassified application routes fail closed.
- Services register annotated controller operations using generated route-group functions; no service implements a second HTTP auth middleware.

A process-local gateway throttle is not a distributed rate limit. Cluster deployments still require a shared edge limiter, a well-defined trusted-proxy configuration, and auth bootstrap. TLS termination and production Origin configuration must be explicit.
