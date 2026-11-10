package request

type ContactValue struct {
	Value     string `json:"value"`
	IsPrimary bool   `json:"isPrimary"`
}

type Contacts struct {
	Emails []ContactValue `json:"emails,omitempty"`
	Phones []ContactValue `json:"phones,omitempty"`
}

// @Validator
// Name fields are normalized and checked before service orchestration.
type CreatePerson struct {
	IdempotencyKey string   `json:"-"`
	FirstName      string   `json:"firstName" normalize:"trim" validate:"notblank,maxrunes=100"`
	MiddleName     string   `json:"middleName,omitempty" normalize:"trim" validate:"maxrunes=100"`
	LastName       string   `json:"lastName" normalize:"trim" validate:"notblank,maxrunes=100"`
	Contacts       Contacts `json:"contacts,omitempty"`
}

// @Validator
// Nil PATCH fields remain absent; supplied fields are validated.
type UpdatePerson struct {
	FirstName  *string   `json:"firstName,omitempty" normalize:"trim" validate:"notblank,maxrunes=100"`
	MiddleName *string   `json:"middleName,omitempty" normalize:"trim" validate:"maxrunes=100"`
	LastName   *string   `json:"lastName,omitempty" normalize:"trim" validate:"notblank,maxrunes=100"`
	Contacts   *Contacts `json:"contacts,omitempty"`
}

type PersonByGID struct {
	GID string `json:"gid"`
}
