package impl

import (
	"context"
	"errors"
	"fmt"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
	authorizationrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/authorization/operations/evaluate"
	"grovehub/services/identity/internal/app/services/security/principal"
)

// @Service
type Service struct {
	repository authorizationrepo.Repository
}

var _ authorization.Service = (*Service)(nil)

func New(repository authorizationrepo.Repository) authorization.Service {
	return &Service{repository: repository}
}

func (service *Service) Require(ctx context.Context, requested permission.Code) error {
	if _, err := principal.RequireTenant(ctx); err != nil {
		return err
	}
	if _, known := permission.Find(requested); !known {
		return authorization.ErrForbidden
	}
	snapshot, err := service.load(ctx)
	if err != nil {
		return err
	}
	if !evaluate.Allowed(ctx, requested, snapshot) {
		return authorization.ErrForbidden
	}
	return nil
}

func (service *Service) Effective(ctx context.Context) ([]permission.Code, error) {
	snapshot, err := service.load(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]permission.Code, 0)
	for _, definition := range permission.Catalog() {
		if evaluate.Allowed(ctx, definition.Code, snapshot) {
			result = append(result, definition.Code)
		}
	}
	return result, nil
}

func (service *Service) load(ctx context.Context) (authorizationrepo.Snapshot, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return authorizationrepo.Snapshot{}, err
	}
	if service.repository == nil {
		return authorizationrepo.Snapshot{}, authorization.ErrUnavailable
	}
	snapshot, err := service.repository.LoadSnapshot(ctx, actor.TenantID, actor.AccountID)
	if errors.Is(err, authorizationrepo.ErrNotFound) {
		return authorizationrepo.Snapshot{}, authorization.ErrForbidden
	}
	if err != nil {
		return authorizationrepo.Snapshot{}, fmt.Errorf("%w: %w", authorization.ErrUnavailable, err)
	}
	if snapshot.AccountID != actor.AccountID || snapshot.TenantID != actor.TenantID || snapshot.Status != entity.Active {
		return authorizationrepo.Snapshot{}, authorization.ErrForbidden
	}
	return snapshot, nil
}
