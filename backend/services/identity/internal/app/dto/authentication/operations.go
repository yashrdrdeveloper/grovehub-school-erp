package authentication
import "grovehub/libs/core/request"
type LoginRequest struct {request.Request;RequestBody Login `validate:"valid"`}
type ForgotPasswordRequest struct {request.Request;RequestBody ForgotPassword `validate:"valid"`}
type ResetPasswordRequest struct {request.Request;RequestBody ResetPassword `validate:"valid"`}
type ChangePasswordRequest struct {request.Request;RequestBody ChangePassword `validate:"valid"`}
type ReauthenticateRequest struct {request.Request;RequestBody Reauthenticate `validate:"valid"`}
type EmptyRequest struct {request.Request}
