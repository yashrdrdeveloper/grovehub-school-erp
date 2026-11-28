package impl

import (
	"context"
	"crypto/ed25519"
	"time"

	coregateway "grovehub/libs/core/gateway"
	"grovehub/libs/platform/scope"
	client "grovehub/services/registry/client/authorization"
	repository "grovehub/services/registry/internal/app/repository/authorization"
	"grovehub/services/registry/internal/app/services/authorization"
	"grovehub/services/registry/internal/app/services/authorization/helpers"
	validator "grovehub/services/registry/internal/app/validators/workload"
)

// @Service
type AuthorizationService struct {
	repository repository.AuthorizationRepository
	audience   string
}

func NewAuthorizationService(repo repository.AuthorizationRepository, audience string) authorization.Service {
	return &AuthorizationService{repository: repo, audience: audience}
}

func (s *AuthorizationService) Verify(ctx context.Context, request client.Request) (context.Context, error) {
	if ctx == nil || ctx.Err() != nil || s.audience == "" || s.repository == nil ||
		scope.ValidateTarget(request.Target) != nil || validator.ValidateIdentifier(request.Capability) != nil ||
		request.Method == "" || request.RequestURI == "" ||
		(request.Target.ProjectID != "" && validator.ValidateID(request.Target.ProjectID) != nil) ||
		(request.Target.TenantID != "" && validator.ValidateID(request.Target.TenantID) != nil) {
		return nil, scope.ErrDenied
	}
	keyID, ok := helpers.KeyID(request.Token)
	if !ok {
		return nil, scope.ErrDenied
	}
	value, err := s.repository.Resolve(ctx, keyID, request.Target, request.Capability, request.AllowInactiveTarget)
	if err != nil || validator.ValidateID(value.WorkloadID) != nil ||
		(value.CallerProjectID != "" && validator.ValidateID(value.CallerProjectID) != nil) ||
		len(value.PublicKey) != ed25519.PublicKeySize {
		return nil, scope.ErrDenied
	}
	// Reuse Core's signature, issuer, subject and expiration verification.
	credentials := coregateway.Ed25519Credentials{Keys: map[string]coregateway.TrustedWorkloadKey{
		keyID: {
			PublicKey: ed25519.PublicKey(value.PublicKey),
			Issuer: value.Issuer,
			Subject: value.Subject,
			Scopes: []string{request.Capability},
		},
	}}
	claims, err := credentials.VerifyCredential(ctx, request.Token)
	now := time.Now()
	if err != nil || claims.Audience != s.audience || len(claims.Scopes) != 1 ||
		claims.Scopes[0] != request.Capability || claims.TargetScope != string(request.Target.Level) ||
		claims.TargetProjectID != request.Target.ProjectID || claims.TargetTenantID != request.Target.TenantID ||
		claims.RequestMethod != request.Method || claims.RequestURI != request.RequestURI ||
		claims.TenantID != "" || claims.ActorGID != "" || claims.SessionID != "" ||
		now.Before(value.ValidFrom) || !now.Before(value.ValidUntil) ||
		claims.IssuedAt.Before(value.ValidFrom) || claims.IssuedAt.Before(value.WorkloadUpdatedAt) ||
		claims.ExpiresAt.After(value.ValidUntil) {
		return nil, scope.ErrDenied
	}
	return scope.WithVerifiedExecution(ctx, scope.Execution{
		Target: request.Target, Capability: request.Capability,
		CallerProjectID: value.CallerProjectID,
		PrincipalType: scope.Workload, PrincipalID: value.WorkloadID, WorkloadID: value.WorkloadID,
	})
}

func (s *AuthorizationService) TenantState(ctx context.Context, projectID, tenantID string) (client.TenantState, error) {
	if ctx == nil || ctx.Err() != nil || s.repository == nil ||
		validator.ValidateID(projectID) != nil || validator.ValidateID(tenantID) != nil {
		return client.TenantState{}, scope.ErrDenied
	}
	value, err := s.repository.TenantState(ctx, projectID, tenantID)
	if err != nil || value.SessionGeneration < 1 {
		return client.TenantState{}, scope.ErrDenied
	}
	return value, nil
}
