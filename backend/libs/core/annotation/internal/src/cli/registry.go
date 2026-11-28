package cli

import (
	"io"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/check"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/clients"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/entities"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/events"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/expand"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/gormtable"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/mappers"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/overlay"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/patchdto"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/routes"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/search"
	testgen "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/testing"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/testing/coverage"
)

func commands() map[string]command {
	return map[string]command{
		"check":     oneArgument("<backend-directory>", check.CheckServices),
		"clients":   oneArgument("<backend-directory>", clients.GenerateAll),
		"entities":  oneArgument("<backend-directory>", entities.GenerateAll),
		"events":    oneArgument("<backend-directory>", events.GenerateAll),
		"expand":    oneArgument("<staging-entity-directory>", expand.Expand),
		"gormtable": twoArguments("<entity-directory> <output-file>", gormtable.GenerateInto),
		"mappers":   oneArgument("<backend-directory>", mappers.GenerateAll),
		"overlay":   twoArguments("<application-directory> <core-directory>", overlay.WriteOverlay),
		"patchdto":  oneArgument("<backend-directory>", patchdto.GenerateAll),
		"routes":    oneArgument("<backend-directory>", routes.GenerateAll),
		"search":    oneArgument("<backend-directory>", search.GenerateAll),
		"testing":   twoArguments("<application-directory> <core-directory>", testgen.GenerateAll),
		"coverage": {
			usage: "[coverage flags]",
			max:   -1,
			run: func(args []string, _ io.Writer) error {
				return coverage.Run(args)
			},
		},
		"inventory": {
			usage: "<package-owner-directory> [--strict]",
			min:   1,
			max:   2,
			run:   runInventory,
		},
	}
}
