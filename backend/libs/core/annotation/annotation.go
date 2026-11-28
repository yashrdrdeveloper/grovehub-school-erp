// Package annotation is the public API for compile-time metadata, source
// scanning, typed DTO rules, presence-aware patches, entity builders and
// startup configuration. Catalog entries alone do not execute runtime policy.
package annotation

import (
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/impl/registry"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/impl/configuration"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/impl/dto"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/impl/entity"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/impl/validation"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/parser"
)

type Status = registry.Status
type Mode = registry.Mode
type Target = registry.Target
type UseCase = registry.UseCase
type Definition = registry.Definition
type Pattern = registry.Pattern
type Directive = registry.Directive
type Occurrence = parser.Occurrence

const (
    Existing = registry.Existing
    Gombok = registry.Gombok
    Proposed = registry.Proposed
    Optional = registry.Optional
    Discouraged = registry.Discouraged
    Pilot = registry.Pilot
    Generate = registry.Generate
    Register = registry.Register
    Check = registry.Check
    Wrap = registry.Wrap
    Tag = registry.Tag
    TypeTarget = registry.TypeTarget
    InterfaceTarget = registry.InterfaceTarget
    MethodTarget = registry.MethodTarget
    FunctionTarget = registry.FunctionTarget
    FieldTarget = registry.FieldTarget
    PackageTarget = registry.PackageTarget
    FileTarget = registry.FileTarget
    ParameterTarget = registry.ParameterTarget
)

var (
    ErrUnknownDirective = registry.ErrUnknownDirective
    ErrInvalidDirective = registry.ErrInvalidDirective
)

func All() []Definition { return registry.All() }
func Lookup(name string) (Definition, bool) { return registry.Lookup(name) }
func Uses() []UseCase { return registry.Uses() }
func Patterns() []Pattern { return registry.Patterns() }
func ParseComment(line string) (Directive, bool, error) { return registry.ParseComment(line) }
func ValidateDirective(d Directive) error { return registry.ValidateDirective(d) }
func ScanFile(path string) ([]Occurrence, error) { return parser.ScanFile(path) }

type Patch[T any] = dto.Patch[T]
type Builder[T any] = entity.Builder[T]

func NewEntity[T any]() Builder[T] { return entity.New[T]() }
func FromEntity[T any](value T) Builder[T] { return entity.From(value) }

type ConfigurationLookup = configuration.Lookup

func ConfigurationFromEnv[T any]() (T, error) { return configuration.FromEnv[T]() }
func DecodeConfiguration[T any](lookup ConfigurationLookup) (T, error) { return configuration.Decode[T](lookup) }

type ValidationError = validation.Error
type ValidationRules[T any] = validation.Rules[T]

func Process[T any](input T, rules ValidationRules[T]) (T, error) { return validation.Process(input, rules) }
func Normalize[T any](input T) (T, error) { return validation.Normalize(input) }
func NormalizeFields[T any](input T) (T, error) { return validation.NormalizeFields(input) }
func Validate[T any](input T) error { return validation.Validate(input) }
