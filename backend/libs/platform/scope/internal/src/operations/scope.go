package operations

import (
	"context"
	"strings"

	"grovehub/libs/platform/scope/internal/src/contracts"
)

type executionKey struct{}

func valid(value string) bool {
	return value != "" && strings.TrimSpace(value) == value && !strings.ContainsRune(value, 0)
}

func ValidateTarget(target contracts.Target) error {
	switch target.Level {
	case contracts.Platform:
		if target.ProjectID == "" && target.TenantID == "" {
			return nil
		}
	case contracts.Project:
		if valid(target.ProjectID) && target.TenantID == "" {
			return nil
		}
	case contracts.Tenant:
		if valid(target.ProjectID) && valid(target.TenantID) {
			return nil
		}
	}
	return contracts.ErrInvalid
}

func ValidateExecution(value contracts.Execution) error {
	if ValidateTarget(value.Target) != nil || !valid(value.PrincipalID) {
		return contracts.ErrInvalid
	}
	if value.CallerProjectID != "" && !valid(value.CallerProjectID) {
		return contracts.ErrInvalid
	}
	switch value.PrincipalType {
	case contracts.User:
		if valid(value.CallerProjectID) && valid(value.SessionID) && valid(value.WorkloadID) {
			return nil
		}
	case contracts.Workload:
		if valid(value.WorkloadID) && value.PrincipalID == value.WorkloadID && value.SessionID == "" {
			return nil
		}
	case contracts.System:
		if value.CallerProjectID == "" && value.SessionID == "" && value.WorkloadID == "" {
			return nil
		}
	}
	return contracts.ErrInvalid
}

func WithVerifiedExecution(ctx context.Context, value contracts.Execution) (context.Context, error) {
	if ctx == nil || ValidateExecution(value) != nil {
		return nil, contracts.ErrInvalid
	}
	return context.WithValue(ctx, executionKey{}, value), nil
}

func From(ctx context.Context) (contracts.Execution, bool) {
	if ctx == nil {
		return contracts.Execution{}, false
	}
	value, ok := ctx.Value(executionKey{}).(contracts.Execution)
	return value, ok
}

func Require(ctx context.Context, target contracts.Target) (contracts.Execution, error) {
	value, ok := From(ctx)
	if !ok || ValidateTarget(target) != nil || value.Target != target {
		return contracts.Execution{}, contracts.ErrDenied
	}
	return value, nil
}

func Allows(grant contracts.Grant, value contracts.Execution, capability string) bool {
	if ValidateExecution(value) != nil || ValidateTarget(grant.Target) != nil ||
		!valid(capability) || !valid(grant.WorkloadID) || grant.Capability != capability ||
		grant.Target != value.Target || grant.CallerProjectID != value.CallerProjectID ||
		grant.WorkloadID != value.WorkloadID {
		return false
	}
	switch value.PrincipalType {
	case contracts.User:
		return grant.AllowDelegation
	case contracts.Workload:
		return true
	default:
		return false
	}
}
