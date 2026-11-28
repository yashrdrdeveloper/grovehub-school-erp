package constants

const Revoke = `UPDATE registry.workload_keys SET revoked_at = now() WHERE workload_id = ? AND key_id = ? AND revoked_at IS NULL`
