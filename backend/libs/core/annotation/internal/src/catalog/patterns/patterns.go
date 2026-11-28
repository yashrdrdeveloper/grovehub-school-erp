// Package patterns owns the 23 GoF pattern definitions and their automation decisions.
// Every entry is metadata; no object is instantiated or decorated automatically.
package patterns

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/contracts"

// All returns a caller-owned copy ordered by GoF family.
func All() []contract.Pattern {
    out := make([]contract.Pattern, 0, len(creational)+len(structural)+len(behavioral))
    out = append(out, creational...)
    out = append(out, structural...)
    return append(out, behavioral...)
}

// Entries is kept for the root annotation registry's historical API.
var Entries = All()
