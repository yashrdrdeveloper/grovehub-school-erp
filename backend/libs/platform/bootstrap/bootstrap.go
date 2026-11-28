// Package bootstrap defines service capabilities consumed by the launcher.
package bootstrap

import (
	"context"

	"grovehub/libs/platform/command"
)

type DatabaseTarget uint8

const (
	BusinessDatabase DatabaseTarget = iota
	AuditDatabase
)

// Migration is a service-owned migration and its required database.
type Migration[D any] struct {
	Target DatabaseTarget
	Run    func(context.Context, D) error
}

// Service is the common registration contract. It does not build HTTP runtimes.
type Service[D any] interface {
	Name() string
	Commands(command.Database[D]) map[string]command.Command
	Migration() Migration[D]
}

// Registration implements Service using service-owned factories.
type Registration[D any] struct {
	ServiceName      string
	RegisterCommands func(command.Database[D]) map[string]command.Command
	MigrationPlan    Migration[D]
}

func (r Registration[D]) Name() string {
	return r.ServiceName
}

func (r Registration[D]) Commands(db command.Database[D]) map[string]command.Command {
	if r.RegisterCommands == nil {
		return nil
	}

	return r.RegisterCommands(db)
}

func (r Registration[D]) Migration() Migration[D] {
	return r.MigrationPlan
}

// RuntimeService extends the common operator contract with typed runtime construction.
// Different services may specialize I and R while sharing Service[D] registration.
type RuntimeService[D any, I any, R any] interface {
	Service[D]
	Build(I) (R, error)
}

// RuntimeRegistration combines operator registration with a service-owned builder.
type RuntimeRegistration[D any, I any, R any] struct {
	Service[D]
	BuildRuntime func(I) (R, error)
}

func (r RuntimeRegistration[D, I, R]) Build(input I) (R, error) {
	return r.BuildRuntime(input)
}
