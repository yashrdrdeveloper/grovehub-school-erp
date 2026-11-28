# Platform Registry

Registry owns technical projects, tenants, workload identities, public keys and
explicit grants. Business records stay in consuming projects. Project codes and
tenant references are data, never compile-time business enumerations.

The first migration creates disabled-by-default projects, tenants and workloads,
with globally unique public keys, exact-scope grants and composite project/tenant
foreign keys. PostgreSQL 15 or newer is required. The CLI registers the migration
through the common bootstrap contract. No API routes or registration writes are
exposed yet. No projects or grants are seeded.

The public entry point is main.go. Annotated entities mirror all six tables.
Project and tenant read services require exact verified scope and capability;
they do not authenticate credentials themselves. Repositories are internal.

Next: provisioning writes, key/grant lifecycle, live credential verification,
annotation-owned restricted APIs and gateway cutover. Runtime authorization is not changed by merely creating these tables.

Project creation requires a platform capability; tenant creation requires a
project capability. Both start disabled. Lifecycle writes require an explicit
resource scope and expected version, and reject stale or no-op transitions.
Disabling a tenant increments its session generation; disabling a project
atomically increments every existing tenant generation. Tenant activation
locks the project and requires it to be active. The gateway must still check
live project and tenant state when it is cut over.

Workload provisioning is deliberately platform-scoped. Registering a workload
does not grant any capabilities. Only canonical Ed25519 public keys are
accepted, with bounded validity. Rotation is add-new then revoke-old, allowing
a controlled overlap; revoked key identifiers and public keys cannot be reused.
Grant creation requires a provider-owned capability and an exact target scope.
Grant revocation deletes the grant. Lifecycle writes require expected versions.
These services are not exposed over HTTP until the gateway is ready.

Live verification is available through Bootstrap().Build(Dependencies) and
Runtime.Authorization(). The route provider must supply an exact target,
capability, method and request URI. The signed credential binds those values
and a deployment-specific audience. The resolver reads active key, caller,
grant and target state in one database statement; no stale-key cache or
legacy fallback is used. The inactive-target exception must be restricted
to provider-defined lifecycle/recovery routes, not request parameters.

The initial verifier accepts workload-only credentials, never actor or
session assertions. Delegation will require independently verified identity
proof. Credentials are short-lived bearer tokens, not single-use nonces;
mutation replay still requires idempotency/expected-version checks. Since
the existing token timestamp has second precision, sign after the next
whole second following key activation or workload status change. No routes
are switched to this verifier in this stage.
