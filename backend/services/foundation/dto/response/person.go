package response

import "time"

type ContactValue struct {
	Value      string `json:"value"`
	IsPrimary  bool   `json:"isPrimary"`
	IsVerified bool   `json:"isVerified"`
}

type Contacts struct {
	Emails []ContactValue `json:"emails"`
	Phones []ContactValue `json:"phones"`
}

type Person struct {
	GID        string    `json:"gid"`
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName,omitempty"`
	LastName   string    `json:"lastName"`
	Contacts   Contacts  `json:"contacts"`
	Status     string    `json:"status"`
	Version    int64     `json:"version"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type PersonErrorMetadata struct {
	GID string `json:"gid,omitempty"`
}
