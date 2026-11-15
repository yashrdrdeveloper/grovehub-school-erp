package response

import "time"

type Connection struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenantId"`
	OwnerGID       string     `json:"ownerGid"`
	Provider       string     `json:"provider"`
	ExternalID     *string    `json:"externalId,omitempty"`
	Status         string     `json:"status"`
	LastTestedAt   *time.Time `json:"lastTestedAt,omitempty"`
	LastTestStatus *string    `json:"lastTestStatus,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
type ConnectionTest struct {
	ID           string    `json:"id"`
	ConnectionID string    `json:"connectionId"`
	Direction    string    `json:"direction"`
	Status       string    `json:"status"`
	Detail       string    `json:"detail"`
	StartedAt    time.Time `json:"startedAt"`
	CompletedAt  time.Time `json:"completedAt"`
}
type OAuthStart struct {
	AuthorizationURL string `json:"authorizationUrl"`
}
type OAuthCompletion struct {
	Status string `json:"status"`
}
