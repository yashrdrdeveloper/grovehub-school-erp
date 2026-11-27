# Restricted OpenAPI migration status

The existing Foundation, Identity and Integration specifications now use restricted Platform paths and workload authentication. Identity and Integration use the dedicated opaque proof header for actor/proof-dependent operations. The Gmail webhook retains its provider-specific JWT requirement.

These are transitional specifications, not complete generated contracts: Identity does not yet describe every recovery/status/RBAC/approval operation; Integration has incomplete response schemas and route inventory; Audit has no complete service specification. The existing generator adapter must not be used to claim complete typed clients until route coverage and PATCH omission/null semantics are verified. Browser-facing School APIs retain their existing public/protected paths independently. No generated client files were produced by this migration.
