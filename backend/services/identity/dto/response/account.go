package response

import "time"

type LoginIdentifier struct {
	Type       string `json:"type"`
	Value      string `json:"value"`
	IsVerified bool   `json:"isVerified"`
}

type BootstrapAccount struct {
	AccountId  string          `json:"accountId"`
	GID        string          `json:"gid"`
	Status     string          `json:"status"`
	Identifier LoginIdentifier `json:"identifier"`
	CreatedAt  time.Time       `json:"createdAt"`
}
