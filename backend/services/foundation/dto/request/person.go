package request

import annotationdto "grovehub/libs/core/annotation/dto"

type ContactValue struct {
	Value     string `json:"value" normalize:"trim" validate:"notblank"`
	IsPrimary bool   `json:"isPrimary"`
}

type Contacts struct {
	Emails []ContactValue `json:"emails,omitempty" validate:"valid"`
	Phones []ContactValue `json:"phones,omitempty" validate:"valid"`
}

// Person groups the common string-valued names for create requests.
type Person struct {
	FirstName  string `json:"firstName" normalize:"trim" validate:"notblank,maxrunes=100"`
	MiddleName string `json:"middleName,omitempty" normalize:"trim" validate:"maxrunes=100"`
	LastName   string `json:"lastName" normalize:"trim" validate:"notblank,maxrunes=100"`
}

type CreatePerson struct {
	Person
	Contacts Contacts `json:"contacts,omitempty" validate:"valid"`
}

// @PatchRequestDto
// UpdatePerson distinguishes absent, null, and supplied values.
// FirstName and LastName cannot be cleared; MiddleName and Contacts can.
type UpdatePerson struct {
	FirstName  annotationdto.Patch[string]   `json:"firstName" normalize:"trim" validate:"nonnull,notblank,maxrunes=100"`
	MiddleName annotationdto.Patch[string]   `json:"middleName" normalize:"trim" validate:"maxrunes=100"`
	LastName   annotationdto.Patch[string]   `json:"lastName" normalize:"trim" validate:"nonnull,notblank,maxrunes=100"`
	Contacts   annotationdto.Patch[Contacts] `json:"contacts" validate:"valid"`
}

type PersonByGID struct {
	GID string `json:"gid"`
}

