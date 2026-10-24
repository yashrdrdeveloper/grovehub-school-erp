package services

import (
	"context"

	identityrequest "grovehub/services/identity/dto/request"
	identityresponse "grovehub/services/identity/dto/response"
)

type RegistrationService interface {
	Register(ctx context.Context, request identityrequest.Register) (identityresponse.Registration, error)
}
