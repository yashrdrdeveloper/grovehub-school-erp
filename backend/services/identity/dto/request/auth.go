package request

import corerequest "grovehub/libs/core/request"

// @Validator
type Register struct {
	FirstName   string  `json:"firstName" normalize:"trim" validate:"notblank,maxrunes=100"`
	// @NotNull
	MiddleName  *string `json:"middleName" normalize:"trim" validate:"nonnull,maxrunes=100"`
	LastName    string  `json:"lastName" normalize:"trim" validate:"notblank,maxrunes=100"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber,omitempty"`
	Password    string  `json:"password"`
}

// RegisterRequest declares only the JSON body, with no unused Params or Query.
type RegisterRequest struct {
	corerequest.Request
	RequestBody Register `validate:"valid"`
}
