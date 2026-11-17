package authentication
type Login struct {Identifier string `json:"identifier"`;Password string `json:"password"`}
type ForgotPassword struct {Email string `json:"email"`}
type ResetPassword struct {NewPassword string `json:"newPassword"`}
type ChangePassword struct {NewPassword string `json:"newPassword"`}
type Reauthenticate struct {Password string `json:"password"`}
