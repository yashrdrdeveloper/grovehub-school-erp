// Package request defines the shared, transport-independent request contract.
package request

import (
	"grovehub/libs/core/request/internal/src/contracts"
	"grovehub/libs/core/request/internal/src/operations/inspect"
)

// Request is the zero-sized marker embedded in concrete endpoint DTOs.
type Request = contracts.Request
type Contract = contracts.Contract
type Sections = contracts.Sections

// Empty and Parts support the original typed REST pipeline.
type Empty = contracts.Empty
type Parts[Path, Query, Body any] = contracts.Parts[Path, Query, Body]

// Inspect discovers only the explicitly declared request sections.
func Inspect(input any) (Sections, error) {
	return inspect.Inspect(input)
}
