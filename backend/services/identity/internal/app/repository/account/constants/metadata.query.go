package constants

const UpdateMetadata = `UPDATE identity.accounts
SET username = ?, search_details = ?, version = version + 1, updated_at = ?
WHERE tenant_id = ? AND account_id = ? AND version = ? AND status = ?
AND status IN ('ACTIVE','PENDING_SETUP')`
