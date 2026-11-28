// Package contract provides compatibility aliases; use annotation's facade instead.
package contract

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"

type Status = contract.Status
type Mode = contract.Mode
type Target = contract.Target
type UseCase = contract.UseCase
type Pattern = contract.Pattern

const (
    Existing = contract.Existing
    Gombok = contract.Gombok
    Proposed = contract.Proposed
    Optional = contract.Optional
    Discouraged = contract.Discouraged
    Pilot = contract.Pilot
    Generate = contract.Generate
    Register = contract.Register
    Check = contract.Check
    Wrap = contract.Wrap
    Tag = contract.Tag
    TypeTarget = contract.TypeTarget
    InterfaceTarget = contract.InterfaceTarget
    MethodTarget = contract.MethodTarget
    FunctionTarget = contract.FunctionTarget
    FieldTarget = contract.FieldTarget
    PackageTarget = contract.PackageTarget
    FileTarget = contract.FileTarget
    ParameterTarget = contract.ParameterTarget
)
