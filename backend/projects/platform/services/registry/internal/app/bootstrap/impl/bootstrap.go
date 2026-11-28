package impl

import (
	"errors"
	"strings"

	"grovehub/services/registry/client/authorization"
	"grovehub/services/registry/internal/app/bootstrap"
	repository "grovehub/services/registry/internal/app/repository/authorization/impl"
	service "grovehub/services/registry/internal/app/services/authorization/impl"
)

type runtime struct {
	authorization authorization.Client
}

func (r *runtime) Authorization() authorization.Client {
	return r.authorization
}

func Build(input bootstrap.Dependencies) (bootstrap.Runtime, error) {
	if input.DB == nil || strings.TrimSpace(input.Audience) == "" || strings.TrimSpace(input.Audience) != input.Audience {
		return nil, errors.New("registry database and deployment audience required")
	}
	return &runtime{authorization: service.NewAuthorizationService(repository.NewAuthorizationRepository(input.DB), input.Audience)}, nil
}
