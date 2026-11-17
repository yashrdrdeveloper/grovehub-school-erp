package account
import "time"
type Account struct {
 AccountID string `json:"accountId"`
 PersonGID string `json:"personGid"`
 Username *string `json:"username"`
 Email string `json:"email"`
 Status string `json:"status"`
 Version int64 `json:"version"`
 CreatedAt time.Time `json:"createdAt"`
 UpdatedAt time.Time `json:"updatedAt"`
}
type CreateResult struct {AccountID string `json:"accountId"`;PersonGID string `json:"personGid"`;Email string `json:"email"`; Status string `json:"status"`;Invitation struct {Status string `json:"status"`} `json:"invitation"`}
type SetupResult struct {AccountID string `json:"accountId"`;PersonGID string `json:"personGid"`;Username string `json:"username"`;Status string `json:"status"`;SetupCompleted bool `json:"setupCompleted"`}
