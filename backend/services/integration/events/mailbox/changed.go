package mailbox

// Changed is a provider notification, not a persisted message.
// Communication will own the mailbox cursor and subsequent ingestion.
type Changed struct {
 Provider string `json:"provider"`
 ConnectionID string `json:"connectionId"`
 TenantID string `json:"tenantId"`
 OwnerGID string `json:"ownerGid"`
 Mailbox string `json:"mailbox"`
 HistoryHint string `json:"historyHint"`
}
const ChangedTopic="integration_mailbox_changed"
const ChangedEventType="integration.mailbox.changed.v1"
