package request

type ContactValue struct {
	Value     string `json:"value" validate:"notblank"`
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

// UpdatePerson distinguishes absent, null, and populated PATCH fields.
// FirstName and LastName cannot be cleared; MiddleName and Contacts can.
type UpdatePerson struct {
	FirstName  PatchString   `json:"firstName"`
	MiddleName PatchString   `json:"middleName"`
	LastName   PatchString   `json:"lastName"`
	Contacts   PatchContacts `json:"contacts"`
}

type PersonByGID struct {
	GID string `json:"gid"`
}
