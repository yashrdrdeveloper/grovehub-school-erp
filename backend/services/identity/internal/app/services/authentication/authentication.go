package authentication

import (
	"context"

	accountdto "grovehub/services/identity/internal/app/dto/account"
	authdto "grovehub/services/identity/internal/app/dto/authentication"
)

type Cookies struct {
	Access string
	Refresh string
}

type SessionResult struct {
	Data authdto.LoginResult
	Cookies Cookies
}

type RefreshResult struct {
	Data authdto.RefreshResult
	Cookies Cookies
}

type Grant struct {
	Secret string
	Data authdto.GrantResult
}

// Service owns the complete authentication use-case surface.
// Login is delegated to the configured credential Strategy. Password setup
// and registration are separate account workflows.
type Service interface {
	Login(context.Context, authdto.Login, string) (SessionResult, error)
	Refresh(context.Context, string) (RefreshResult, error)
	Forgot(context.Context, string) error
	Exchange(context.Context, string, string) (Grant, error)
	ResetPassword(context.Context, string, string) error
	Reauthenticate(context.Context, string) error
	ChangePassword(context.Context, string) error
	Logout(context.Context) error
	LogoutAll(context.Context) (int64, error)
	ListSessions(context.Context) (authdto.SessionsResult, error)
	RevokeSession(context.Context, string) error
	Me(context.Context) (accountdto.Account, error)
}
