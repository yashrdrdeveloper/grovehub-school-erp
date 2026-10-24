package response

import "time"

type Domain struct {
	DomainId          string    `json:"domainId"`
	Name              string    `json:"name"`
	VerificationHost  string    `json:"verificationHost"`
	VerificationValue string    `json:"verificationValue"`
	Status            string    `json:"status"`
	Version           int64     `json:"version"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
