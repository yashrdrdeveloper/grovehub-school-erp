package constants

const FindByID = "workload_id = ?"

const ChangeStatus = `UPDATE registry.workloads
SET status = ?, version = version + 1, updated_at = now()
WHERE workload_id = ? AND version = ? AND status <> ?
RETURNING *`
