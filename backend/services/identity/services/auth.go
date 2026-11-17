package services

import (
	"context"

	identityrequest "grovehub/services/identity/dto/request"
	identityresponse "grovehub/services/identity/dto/response"
)

type AuthService interface {
	Register(ctx context.Context, request identityrequest.RegisterRequest) (identityresponse.Register, error)
	Login(ctx context.Context, request identityrequest.LoginRequest) (identityresponse.Login, error)
	SetPassword(ctx context.Context, request identityrequest.SetPasswordRequest) (identityresponse.SetPassword, error)
}
