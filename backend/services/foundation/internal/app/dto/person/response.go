package person

import "time"

type ResponseContactValue struct {
	Value      string `json:"value"`
	IsPrimary  bool   `json:"isPrimary"`
	IsVerified bool   `json:"isVerified"`
}

type ResponseContacts struct {
	Emails []ResponseContactValue `json:"emails"`
	Phones []ResponseContactValue `json:"phones"`
}

type Person struct {
	GID        string    `json:"gid"`
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName,omitempty"`
	LastName   string    `json:"lastName"`
	Contacts   ResponseContacts  `json:"contacts"`
	Status     string    `json:"status"`
	Version    int64     `json:"version"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type PersonErrorMetadata struct {
	GID string `json:"gid,omitempty"`
}
