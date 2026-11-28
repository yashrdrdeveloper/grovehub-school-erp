// Package dto declares build-time DTO annotations and shared PATCH support.
package dto

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"

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
