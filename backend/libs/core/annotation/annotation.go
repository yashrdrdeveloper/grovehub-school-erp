// Package annotation exposes metadata, source scanning, typed DTO rules,
// presence-aware patches, entity builders and startup configuration.
package annotation

import (
	"grovehub/libs/core/annotation/internal/src/configuration"
	"grovehub/libs/core/annotation/internal/src/contracts"
	"grovehub/libs/core/annotation/internal/src/dto"
	"grovehub/libs/core/annotation/internal/src/entity"
	"grovehub/libs/core/annotation/internal/src/parser"
	"grovehub/libs/core/annotation/internal/src/registry"
	"grovehub/libs/core/annotation/internal/src/validation"
)

type (
	Status     = contract.Status
	Mode       = contract.Mode
	Target     = contract.Target
	UseCase    = contract.UseCase
	Definition = contract.Definition
	Pattern    = contract.Pattern
	Directive  = contract.Directive
	Occurrence = contract.Occurrence
)

const (
	Existing    = contract.Existing
	Gombok      = contract.Gombok
	Proposed    = contract.Proposed
	Optional    = contract.Optional
	Discouraged = contract.Discouraged
	Pilot       = contract.Pilot
)

const (
	Generate = contract.Generate
	Register = contract.Register
	Check    = contract.Check
	Wrap     = contract.Wrap
	Tag      = contract.Tag
)

const (
	TypeTarget      = contract.TypeTarget
	InterfaceTarget = contract.InterfaceTarget
	MethodTarget    = contract.MethodTarget
	FunctionTarget  = contract.FunctionTarget
	FieldTarget     = contract.FieldTarget
	PackageTarget   = contract.PackageTarget
	FileTarget      = contract.FileTarget
	ParameterTarget = contract.ParameterTarget
)

var (
	ErrUnknownDirective = registry.ErrUnknownDirective
	ErrInvalidDirective = registry.ErrInvalidDirective
)

func All() []Definition {
	return registry.All()
}

func Lookup(name string) (Definition, bool) {
	return registry.Lookup(name)
}

func Uses() []UseCase {
	return registry.Uses()
}

func Patterns() []Pattern {
	return registry.Patterns()
}

func ParseComment(line string) (Directive, bool, error) {
	return registry.ParseComment(line)
}

func ValidateDirective(d Directive) error {
	return registry.ValidateDirective(d)
}

func ScanFile(path string) ([]Occurrence, error) {
	return parser.ScanFile(path)
}

type Patch[T any] = dto.Patch[T]

type Builder[T any] = entity.Builder[T]

func NewEntity[T any]() Builder[T] {
	return entity.New[T]()
}

func FromEntity[T any](value T) Builder[T] {
	return entity.From(value)
}

type ConfigurationLookup = configuration.Lookup

func ConfigurationFromEnv[T any]() (T, error) {
	return configuration.FromEnv[T]()
}

func DecodeConfiguration[T any](lookup ConfigurationLookup) (T, error) {
	return configuration.Decode[T](lookup)
}

type ValidationError = validation.Error

type ValidationRules[T any] = validation.Rules[T]

func Process[T any](input T, rules ValidationRules[T]) (T, error) {
	return validation.Process(input, rules)
}

func Normalize[T any](input T) (T, error) {
	return validation.Normalize(input)
}

func NormalizeFields[T any](input T) (T, error) {
	return validation.NormalizeFields(input)
}

func Validate[T any](input T) error {
	return validation.Validate(input)
}
