// Package dto declares build-time DTO annotations and shared PATCH support.
package dto

import "grovehub/libs/core/annotation/contract"

// UseCases registers annotations implemented by the DTO generator.
var UseCases = []contract.UseCase{
	{
		Name: "PatchRequestDto",
		Syntax: "@PatchRequestDto",
		Category: "DTO request serialization",
		Purpose: "Generate presence-aware JSON MarshalJSON for fields typed dto.Patch[T]; omit absent fields, retain explicit null",
		Mode: "GEN",
		Status: contract.Pilot,
		Roadmap: "P2",
		Target: contract.TypeTarget,
	},
}
