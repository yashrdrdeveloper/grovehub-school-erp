package changed
// Change represents a safe account event snapshot with no credential material.
type Snapshot struct{
 AccountID string `json:"accountId"`
 PersonGID string `json:"personGid"`
 TenantID string `json:"tenantId"`
 Username *string `json:"username,omitempty"`
 Status string `json:"status"`
 Version int64 `json:"version"`
}
// @Event topic=identity_account_changed type=identity.account.changed.v1 source=identity
type Changed struct{
 Operation string `json:"operation"`
 Old *Snapshot `json:"old,omitempty"`
 New *Snapshot `json:"new,omitempty"`
}
