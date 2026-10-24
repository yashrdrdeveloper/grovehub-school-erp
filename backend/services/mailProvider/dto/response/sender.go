package response

import "time"

type Sender struct {
	SenderId    string    `json:"senderId"`
	DomainId    string    `json:"domainId"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	Status      string    `json:"status"`
	Version     int64     `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
