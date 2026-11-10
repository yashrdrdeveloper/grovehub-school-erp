// Package patterns owns the 23 GoF pattern definitions and their automation decisions.
// Every entry is metadata; no object is instantiated or decorated automatically.
package patterns

import "grovehub/libs/core/annotation/contract"

// All returns a caller-owned copy ordered by GoF family.
func All() []contract.Pattern {
    out := make([]contract.Pattern, 0, len(creational)+len(structural)+len(behavioral))
    out = append(out, creational...)
    out = append(out, structural...)
    return append(out, behavioral...)
}

// Entries is kept for the root annotation registry's historical API.
var Entries = All()
