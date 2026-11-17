package authentication

import "time"

type LoginResult struct {
	AccountID     string `json:"accountId"`
	PersonGID     string `json:"personGid"`
	SessionID     string `json:"sessionId"`
	ExpiresIn     int    `json:"expiresIn"`
	RequiresSetup bool   `json:"requiresSetup"`
}

type RefreshResult struct {
	SessionID string `json:"sessionId"`
	ExpiresIn int    `json:"expiresIn"`
	Refreshed bool   `json:"refreshed"`
}

type StatusResult struct {
	Status          string `json:"status"`
	SessionsRevoked bool   `json:"sessionsRevoked,omitempty"`
}

type AcceptedResult struct {
	Accepted bool   `json:"accepted"`
	Delivery string `json:"delivery"`
}

type GrantResult struct {
	Verified  bool   `json:"verified"`
	Purpose   string `json:"purpose"`
	ExpiresIn int    `json:"expiresIn"`
}

type SessionSummary struct {
	SessionID string `json:"sessionId"`
	Current bool `json:"current"`
	CreatedAt time.Time `json:"createdAt"`
	LastActiveAt *time.Time `json:"lastActiveAt,omitempty"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type SessionsResult struct {
	Items []SessionSummary `json:"items"`
}

type RevocationResult struct {
	Status string `json:"status"`
	RevokedCount int64 `json:"revokedCount"`
}

type SessionRevokedResult struct {
	SessionID string `json:"sessionId"`
	Status string `json:"status"`
}
