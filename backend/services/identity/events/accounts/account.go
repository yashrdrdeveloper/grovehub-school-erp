package accounts

const (
	CreatedTopic = "identity_account_created"
	CreatedType  = "identity.account.created.v1"
	Source       = "identity"
)

type Created struct {
	AccountID string `json:"accountId"`
	PersonGID string `json:"personGid"`
	Status    string `json:"status"`
	Version   int64  `json:"version"`
}
