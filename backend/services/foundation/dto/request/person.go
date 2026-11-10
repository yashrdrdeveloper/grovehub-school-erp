package request

type ContactValue struct {
	// @NotBlank
	Value string `json:"value" validate:"notblank"`
	IsPrimary bool `json:"isPrimary"`
}

type Contacts struct {
	// @Valid
	Emails []ContactValue `json:"emails,omitempty" validate:"valid"`
	// @Valid
	Phones []ContactValue `json:"phones,omitempty" validate:"valid"`
}

// @Data
// Person is the shared name contract; string is used for create, *string for PATCH.
type Person[T any] struct {
	// @NotBlank
	// @Size max=100
	FirstName T `json:"firstName,omitempty" normalize:"trim" validate:"notblank,maxrunes=100"`
	// @Size max=100
	MiddleName T `json:"middleName,omitempty" normalize:"trim" validate:"maxrunes=100"`
	// @NotBlank
	// @Size max=100
	LastName T `json:"lastName,omitempty" normalize:"trim" validate:"notblank,maxrunes=100"`
}

// @Validator
// Header metadata is not part of the JSON request body.
type CreatePerson struct {
	// @Valid
	Person[string] `validate:"valid"`
	// @Valid
	Contacts Contacts `json:"contacts,omitempty" validate:"valid"`
}

// @Validator
// Nil PATCH fields remain absent; supplied values are validated.
type UpdatePerson struct {
	// @Valid
	Person[*string] `validate:"valid"`
	// @Valid
	Contacts *Contacts `json:"contacts,omitempty" validate:"valid"`
}

type PersonByGID struct {
	GID string `json:"gid"`
}
