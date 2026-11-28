// Package scope defines explicit, authenticated execution boundaries.
// It has no dependency on transport, persistence or consuming projects.
package scope

import (
	"context"

	"grovehub/libs/platform/scope/internal/src/contracts"
	"grovehub/libs/platform/scope/internal/src/operations"
)

type Level = contracts.Level
type PrincipalType = contracts.PrincipalType
type Target = contracts.Target
type Execution = contracts.Execution
type Grant = contracts.Grant

const (
	Platform = contracts.Platform
	Project = contracts.Project
	Tenant = contracts.Tenant
	User = contracts.User
	Workload = contracts.Workload
	System = contracts.System
)

var ErrInvalid = contracts.ErrInvalid
var ErrDenied = contracts.ErrDenied

func ValidateTarget(target Target) error {
	return operations.ValidateTarget(target)
}

func ValidateExecution(value Execution) error {
	return operations.ValidateExecution(value)
}

// WithVerifiedExecution is for trusted authentication adapters only.
// Structural validation is not credential verification.
func WithVerifiedExecution(ctx context.Context, value Execution) (context.Context, error) {
	return operations.WithVerifiedExecution(ctx, value)
}

func From(ctx context.Context) (Execution, bool) {
	return operations.From(ctx)
}

// Require denies absent contexts and nonmatching targets. A project grant
// never implicitly authorizes a tenant operation.
func Require(ctx context.Context, target Target) (Execution, error) {
	return operations.Require(ctx, target)
}

// Allows checks one already verified grant; it never grants broader scope.
func Allows(grant Grant, value Execution, capability string) bool {
	return operations.Allows(grant, value, capability)
}
