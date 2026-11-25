package account

import "time"

type Account struct {
	AccountID string    `json:"accountId"`
	PersonGID string    `json:"personGid"`
	Username  *string   `json:"username"`
	Profile *Profile `json:"profile,omitempty"`
	PrimaryRoleID *string `json:"primaryRoleId"`
	SecondaryRoleIDs []string `json:"secondaryRoleIds"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	Version   int64     `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type InvitationStatus struct {
	Status string `json:"status"`
}

type CreateResult struct {
	AccountID  string           `json:"accountId"`
	PersonGID  string           `json:"personGid"`
	Email      string           `json:"email"`
	Status     string           `json:"status"`
	Invitation InvitationStatus `json:"invitation"`
}

type SetupResult struct {
	AccountID      string `json:"accountId"`
	PersonGID      string `json:"personGid"`
	Status         string `json:"status"`
	SetupCompleted bool   `json:"setupCompleted"`
}

type SetupStatusResult struct {
	Verified bool `json:"verified"`
	ExpiresIn int `json:"expiresIn"`
}
