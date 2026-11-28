# Foundation authorization boundary

Foundation persons are global identities, not tenant-owned records. Do not expose unrestricted person search or mutation to School browser sessions. Only explicitly configured school-foundation service workload keys may call the five restricted person operations, with distinct scopes. A future School membership/use-case service must establish tenant-specific authorization before exposing browser wrappers. Do not reuse the ordinary School browser-wrapper signing key for Foundation grants.
