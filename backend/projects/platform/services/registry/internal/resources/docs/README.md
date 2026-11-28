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
