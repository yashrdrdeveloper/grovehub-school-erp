package impl

import (
	"context"
	"errors"
	"time"

	coreid "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/id"
	kind "grovehub/services/identity/internal/app/constants/approval"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	repo "grovehub/services/identity/internal/app/repository/approval"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	authrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/repository/authority"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/services/approval"
	"grovehub/services/identity/internal/app/services/approval/helpers"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/principal"
)

// @Service
type Service struct {
	repository repo.Repository
	locker authority.Locker
	authorizer authorization.Service
	snapshots authrepo.Repository
	accounts accountrepo.Reader
	roles rolerepo.Repository
}

var _ approval.Service = (*Service)(nil)

func New(repository repo.Repository, locker authority.Locker, authorizer authorization.Service, snapshots authrepo.Repository, accounts accountrepo.Reader, roles rolerepo.Repository) approval.Service {
	return &Service{repository: repository, locker: locker, authorizer: authorizer, snapshots: snapshots, accounts: accounts, roles: roles}
}

func (service *Service) available() bool {
	return service.repository != nil && service.locker != nil && service.authorizer != nil && service.snapshots != nil && service.accounts != nil && service.roles != nil
}

func (service *Service) CreateLifecycle(ctx context.Context, request dto.CreateLifecycle) (entity.ApprovalRequest, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	if !service.available() {
		return entity.ApprovalRequest{}, authorization.ErrUnavailable
	}
	request, code, payload, digest, err := helpers.NormalizeLifecycle(request)
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	var created entity.ApprovalRequest
	err = service.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		if err := service.authorizer.Require(tx, code); err != nil {
			return err
		}
		existing, err := service.repository.FindByKey(tx, actor.TenantID, actor.AccountID, request.IdempotencyKey)
		if err == nil {
			if existing.RequestType != request.Type || existing.PayloadSHA256 != digest || existing.Reason != request.Reason {
				return repo.ErrConflict
			}
			created = existing
			return nil
		}
		if !errors.Is(err, repo.ErrNotFound) {
			return err
		}
		if err := service.validateTarget(tx, actor.TenantID, request); err != nil {
			return err
		}
		policy, err := service.repository.ActivePolicy(tx, actor.TenantID, request.Type)
		if err != nil {
			return err
		}
		if err := helpers.ValidatePolicy(policy); err != nil {
			return err
		}
		now, err := service.repository.Now(tx)
		if err != nil {
			return err
		}
		id, err := coreid.NewUUIDv7()
		if err != nil {
			return err
		}
		created = entity.ApprovalRequest{RequestID: id, TenantID: actor.TenantID, PolicyID: policy.PolicyID, RequesterID: actor.AccountID, RequestType: request.Type, Payload: payload, PayloadSHA256: digest, Reason: request.Reason, IdempotencyKey: request.IdempotencyKey, Status: kind.Pending, CurrentStep: 1, StepCount: len(policy.Steps), Version: 1, ExpiresAt: now.Add(24*time.Hour), CreatedAt: now, UpdatedAt: now}
		return service.repository.Insert(tx, created, policy.Steps)
	})
	if err != nil {
		return entity.ApprovalRequest{}, err
	}
	return created, nil
}
