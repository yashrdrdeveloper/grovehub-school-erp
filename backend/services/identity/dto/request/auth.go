package request

// @Validator
type Register struct {

    FirstName      string  `json:"firstName" validate:"notblank"`
    // @NotNull
    MiddleName *string `json:"middleName" validate:"nonnull"`
    LastName       string  `json:"lastName" validate:"notblank"`
    Email          string  `json:"email"`
    PhoneNumber    string  `json:"phoneNumber,omitempty"`
    Password       string  `json:"password"`
}
