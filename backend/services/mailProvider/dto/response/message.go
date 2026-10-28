package response

import "time"

type Message struct {
	MessageId      string    `json:"messageId"`
	GID            string    `json:"gid"`
	SenderId       string    `json:"senderId"`
	To             []string  `json:"to"`
	Cc             []string  `json:"cc,omitempty"`
	Bcc            []string  `json:"bcc,omitempty"`
	Subject        string    `json:"subject"`
	Status         string    `json:"status"`
	FailureMessage string    `json:"failureMessage,omitempty"`
	Version        int64     `json:"version"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
