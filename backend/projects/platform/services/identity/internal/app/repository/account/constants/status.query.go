package constants

const ChangeStatus = `UPDATE identity.accounts
SET status = ?, version = version + 1, updated_at = ?
WHERE tenant_id = ? AND account_id = ? AND version = ? AND status = ?`
