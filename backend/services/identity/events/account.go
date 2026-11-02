package events

const (
	AccountCreatedTopic = "identity_account_created"
	AccountCreatedType  = "identity.account.created.v1"
	Source              = "identity"
)

type AccountCreated struct {
	AccountID string `json:"accountId"`
	PersonGID string `json:"personGid"`
	Status    string `json:"status"`
	Version   int64  `json:"version"`
}
