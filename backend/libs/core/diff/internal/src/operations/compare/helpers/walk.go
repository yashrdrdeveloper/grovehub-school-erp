package helpers

import (
	"reflect"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/diff/internal/src/contracts"
)

// Changes walks JSON-normalized values and returns deterministic leaf changes.
func Changes(old, next any) []contracts.Change {
	changes := make([]contracts.Change, 0)
	walk("", old, next, &changes)

	return changes
}

func walk(path string, old, next any, changes *[]contracts.Change) {
	if reflect.DeepEqual(old, next) {
		return
	}

	oldMap, oldIsMap := old.(map[string]any)
	newMap, newIsMap := next.(map[string]any)

	if (oldIsMap && newIsMap) || (oldIsMap && next == nil) || (newIsMap && old == nil) {
		walkObject(path, oldMap, newMap, changes)
		return
	}

	oldList, oldIsList := old.([]any)
	newList, newIsList := next.([]any)

	if (oldIsList && newIsList) || (oldIsList && next == nil) || (newIsList && old == nil) {
		walkArray(path, oldList, newList, changes)
		return
	}

	if path == "" {
		path = "$"
	}

	*changes = append(*changes, contracts.Change{
		Field:    path,
		OldValue: old,
		NewValue: next,
	})
}
