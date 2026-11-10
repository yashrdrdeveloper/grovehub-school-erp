package accounts

// @Event topic=identity_account_created type=identity.account.created.v1 source=identity
type Created struct {
	AccountID string `json:"accountId"`
	PersonGID string `json:"personGid"`
	Status    string `json:"status"`
	Version   int64  `json:"version"`
}
