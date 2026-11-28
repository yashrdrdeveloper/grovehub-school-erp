package helpers

import (
	"fmt"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/registry"
)

func ValidateOccurrence(found contract.Occurrence) error {
	if found.Directive.Name == "Provider" {
		if _, _, err := ParseProviderDirective(found.Directive.Arguments); err != nil {
			return fmt.Errorf("%s:%d: %w", found.File, found.Line, err)
		}
	}

	definition, exists := registry.Lookup(found.Directive.Name)
	if !exists {
		return fmt.Errorf("%s:%d: unknown annotation @%s", found.File, found.Line, found.Directive.Name)
	}

	return validateTarget(found, definition)
}

func validateTarget(found contract.Occurrence, definition contract.Definition) error {
	matchedTarget := false

	for _, useCase := range definition.UseCases {
		if useCase.Target != found.Target {
			continue
		}

		matchedTarget = true

		if isExecutable(useCase.Status) {
			return nil
		}
	}

	if matchedTarget {
		return fmt.Errorf("%s:%d: @%s is catalog metadata only, not an enforced processor; remove it until implemented", found.File, found.Line, found.Directive.Name)
	}

	return fmt.Errorf("%s:%d: @%s cannot annotate target %s", found.File, found.Line, found.Directive.Name, found.Target)
}

func isExecutable(status contract.Status) bool {
	return status == contract.Existing || status == contract.Pilot
}
