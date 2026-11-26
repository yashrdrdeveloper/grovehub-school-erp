package constants

// Global authority must be acquired before any tenant authority lock.
// Shared locks permit independent tenants to proceed concurrently.
const LockGlobalShared = `SELECT pg_advisory_xact_lock_shared(hashtextextended('identity:authority:global:v1', 0))`
const LockGlobalExclusive = `SELECT pg_advisory_xact_lock(hashtextextended('identity:authority:global:v1', 0))`

const Tenants = `SELECT tenant_id FROM identity.accounts WHERE tenant_id IS NOT NULL AND btrim(tenant_id) <> ''
UNION SELECT tenant_id FROM identity.rbac_migration_receipts WHERE operation = 'BOOTSTRAP'
ORDER BY tenant_id`

const Bootstrapped = `SELECT EXISTS (SELECT 1 FROM identity.rbac_migration_receipts WHERE tenant_id = ? AND operation = 'BOOTSTRAP')`
