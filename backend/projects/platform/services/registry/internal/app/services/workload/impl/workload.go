package impl

import (
	"context"
	"time"

	"grovehub/libs/core/id"
	"grovehub/libs/core/utils"
	"grovehub/services/registry/internal/app/constants/capability"
	"grovehub/services/registry/internal/app/constants/status"
	"grovehub/services/registry/internal/app/constants/validation"
	dto "grovehub/services/registry/internal/app/dto/workload"
	"grovehub/services/registry/internal/app/entity"
	"grovehub/services/registry/internal/app/repository/grant"
	"grovehub/services/registry/internal/app/repository/key"
	repository "grovehub/services/registry/internal/app/repository/workload"
	"grovehub/services/registry/internal/app/services/helpers"
	"grovehub/services/registry/internal/app/services/workload"
	access "grovehub/services/registry/internal/app/services/workload/helpers"
	validator "grovehub/services/registry/internal/app/validators/workload"
)

// @Service
type WorkloadService struct {
	workloads repository.WorkloadRepository
	keys      key.KeyRepository
	grants    grant.GrantRepository
}

func NewWorkloadService(workloads repository.WorkloadRepository, keys key.KeyRepository, grants grant.GrantRepository) workload.WorkloadService {
	return &WorkloadService{workloads: workloads, keys: keys, grants: grants}
}

func (s *WorkloadService) Get(ctx context.Context, workloadID string) (entity.Workload, error) {
	if err := access.Require(ctx, capability.WorkloadRead); err != nil {
		return entity.Workload{}, err
	}
	if err := validator.ValidateID(workloadID); err != nil {
		return entity.Workload{}, err
	}
	value, err := s.workloads.Find(ctx, workloadID)
	return value, helpers.RepositoryError(err)
}

func (s *WorkloadService) Create(ctx context.Context, input dto.CreateInput) (entity.Workload, error) {
	if err := access.Require(ctx, capability.WorkloadCreate); err != nil {
		return entity.Workload{}, err
	}
	if err := validator.ValidateCreate(input); err != nil {
		return entity.Workload{}, err
	}
	gid, err := id.NewUUIDv7()
	if err != nil {
		return entity.Workload{}, helpers.RepositoryError(err)
	}
	now := time.Now().UTC()
	value, err := s.workloads.Create(ctx, entity.Workload{
		WorkloadID: gid,
		ProjectID: input.ProjectID,
		Issuer: input.Issuer,
		Subject: input.Subject,
		Status: status.Disabled,
		Version: validation.InitialVersion,
		CreatedAt: now,
		UpdatedAt: now,
	})
	return value, helpers.RepositoryError(err)
}

func (s *WorkloadService) ChangeStatus(ctx context.Context, workloadID string, input dto.StatusInput) (entity.Workload, error) {
	if err := access.Require(ctx, capability.WorkloadStatus); err != nil {
		return entity.Workload{}, err
	}
	if err := validator.ValidateID(workloadID); err != nil {
		return entity.Workload{}, err
	}
	if err := validator.ValidateStatus(input); err != nil {
		return entity.Workload{}, err
	}
	value, err := s.workloads.ChangeStatus(ctx, workloadID, input.ExpectedVersion, input.Status)
	return value, helpers.RepositoryError(err)
}

func (s *WorkloadService) RegisterKey(ctx context.Context, workloadID string, input dto.KeyInput) error {
	if err := access.Require(ctx, capability.KeyRegister); err != nil {
		return err
	}
	if err := validator.ValidateID(workloadID); err != nil {
		return err
	}
	now := time.Now().UTC()
	publicKey, err := validator.ValidateKey(input, now)
	if err != nil {
		return err
	}
	return helpers.RepositoryError(s.keys.Create(ctx, entity.WorkloadKey{
		KeyID: input.KeyID,
		WorkloadID: workloadID,
		PublicKey: publicKey,
		ValidFrom: input.ValidFrom.UTC(),
		ValidUntil: input.ValidUntil.UTC(),
		CreatedAt: now,
	}))
}

func (s *WorkloadService) RevokeKey(ctx context.Context, workloadID, keyID string) error {
	if err := access.Require(ctx, capability.KeyRevoke); err != nil {
		return err
	}
	if err := validator.ValidateID(workloadID); err != nil {
		return err
	}
	if err := validator.ValidateIdentifier(keyID); err != nil {
		return err
	}
	return helpers.RepositoryError(s.keys.Revoke(ctx, workloadID, keyID))
}

func (s *WorkloadService) CreateGrant(ctx context.Context, workloadID string, input dto.GrantInput) (string, error) {
	if err := access.Require(ctx, capability.GrantCreate); err != nil {
		return "", err
	}
	if err := validator.ValidateID(workloadID); err != nil {
		return "", err
	}
	if err := validator.ValidateGrant(input); err != nil {
		return "", err
	}
	gid, err := id.NewUUIDv7()
	if err != nil {
		return "", helpers.RepositoryError(err)
	}
	err = s.grants.Create(ctx, entity.WorkloadGrant{
		GrantID: gid,
		WorkloadID: workloadID,
		TargetProjectID: utils.OptionalString(input.TargetProjectID),
		TargetTenantID: utils.OptionalString(input.TargetTenantID),
		Scope: input.Scope,
		Capability: input.Capability,
		AllowDelegation: input.AllowDelegation,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return "", helpers.RepositoryError(err)
	}
	return gid, nil
}

func (s *WorkloadService) RevokeGrant(ctx context.Context, workloadID, grantID string) error {
	if err := access.Require(ctx, capability.GrantRevoke); err != nil {
		return err
	}
	if err := validator.ValidateID(workloadID); err != nil {
		return err
	}
	if err := validator.ValidateID(grantID); err != nil {
		return err
	}
	return helpers.RepositoryError(s.grants.Revoke(ctx, workloadID, grantID))
}
