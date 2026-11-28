# Platform Registry

Registry owns technical projects, tenants, workload identities, public keys and
explicit grants. Business records stay in consuming projects. Project codes and
tenant references are data, never compile-time business enumerations.

The first migration creates disabled-by-default projects, tenants and workloads,
with globally unique public keys, exact-scope grants and composite project/tenant
foreign keys. PostgreSQL 15 or newer is required. The CLI registers the migration
through the common bootstrap contract. No API routes or registration writes are
exposed yet. No projects or grants are seeded.

Next: repository and service interfaces, trusted operator provisioning,
annotation-owned restricted APIs, live credential verification and gateway
cutover. Runtime authorization is not changed by merely creating these tables.
