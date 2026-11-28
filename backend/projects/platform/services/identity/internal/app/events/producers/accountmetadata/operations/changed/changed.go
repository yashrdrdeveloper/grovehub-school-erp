package changed

type Snapshot struct {
	AccountID string `json:"accountId"`
	Username *string `json:"username"`
	Version int64 `json:"version"`
}

type Changed struct {
	TenantID string `json:"tenantId"`
	ActorID string `json:"actorId"`
	Operation string `json:"operation"`
	Old Snapshot `json:"old"`
	New Snapshot `json:"new"`
}
