package mailbox

import coreevents "grovehub/libs/core/events"

// Payload describes the provider notification only.
// Communication owns message state and the authoritative sync cursor.
type Payload struct {
	Provider     string `json:"provider"`
	ConnectionID string `json:"connectionId"`
	TenantID     string `json:"tenantId"`
	OwnerGID     string `json:"ownerGid"`
	Mailbox      string `json:"mailbox"`
	HistoryHint  string `json:"historyHint"`
}

// Changed uses the shared generic event envelope, never a custom event wrapper.
// @Event topic=integration_mailbox_changed type=integration.mailbox.changed.v1 typeConst=ChangedEventType
type Changed = coreevents.Envelope[Payload]

