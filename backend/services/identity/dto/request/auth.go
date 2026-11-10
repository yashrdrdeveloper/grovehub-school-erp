package request

// @Validator
type Register struct {
    IdempotencyKey string  `json:"-"`
    FirstName      string  `json:"firstName" validate:"notblank"`
    MiddleName     *string `json:"middleName" validate:"required"`
    LastName       string  `json:"lastName" validate:"notblank"`
    Email          string  `json:"email"`
    PhoneNumber    string  `json:"phoneNumber,omitempty"`
    Password       string  `json:"password"`
}
