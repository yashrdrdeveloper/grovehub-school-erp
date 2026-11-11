package request

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
