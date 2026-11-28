# Platform gateway session boundary

Import `grovehub/launcher/services/gateway/session` for authenticated principals, cookie constants, access verification, protected route authorization and tenant-generation checks. The session facade never imports the gateway runtime, Identity or Integration; both runtime and services use the same context keys and contract types. Authentication and tenant checks retain their existing semantics. Do not accept principals or trusted tenant checks from HTTP input. The former library is only a compatibility bridge during migration.
