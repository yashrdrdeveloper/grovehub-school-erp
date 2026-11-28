package flows

type Payload struct {
	Provider     string `json:"provider"`
	ConnectionID string `json:"connectionId"`
	TenantID     string `json:"tenantId"`
	OwnerGID     string `json:"ownerGid"`
	Mailbox      string `json:"mailbox"`
	HistoryHint  string `json:"historyHint"`
}
