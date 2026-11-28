package constants

const Revoke = `DELETE FROM registry.workload_grants WHERE workload_id = ? AND grant_id = ?`
