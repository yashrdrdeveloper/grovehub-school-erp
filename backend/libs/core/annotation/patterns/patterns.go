// Package patterns provides a compatibility API for GoF metadata.
package patterns

import (
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"
    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/catalog/patterns"
)

func All() []contract.Pattern { return patterns.All() }

var Entries = patterns.Entries
