package authentication
type LoginResult struct {AccountID string `json:"accountId"`;PersonGID string `json:"personGid"`;SessionID string `json:"sessionId"`; ExpiresIn int `json:"expiresIn"`;RequiresSetup bool `json:"requiresSetup"`}
type RefreshResult struct {SessionID string `json:"sessionId"`;ExpiresIn int `json:"expiresIn"`;Refreshed bool `json:"refreshed"`}
type StatusResult struct {Status string `json:"status"`;SessionsRevoked bool `json:"sessionsRevoked,omitempty"`}
type AcceptedResult struct {Accepted bool `json:"accepted"`;Delivery string `json:"delivery"`}
type GrantResult struct {Verified bool `json:"verified"`;Purpose string `json:"purpose"`;ExpiresIn int `json:"expiresIn"`}
