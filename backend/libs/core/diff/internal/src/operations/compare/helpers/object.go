package helpers

import (
	"sort"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/diff/internal/src/contracts"
)

func walkObject(path string, old, next map[string]any, changes *[]contracts.Change) {
	for _, name := range objectKeys(old, next) {
		walk(objectPath(path, name), old[name], next[name], changes)
	}
}

func objectKeys(old, next map[string]any) []string {
	names := make(map[string]struct{}, len(old)+len(next))

	for name := range old {
		names[name] = struct{}{}
	}

	for name := range next {
		names[name] = struct{}{}
	}

	keys := make([]string, 0, len(names))
	for name := range names {
		keys = append(keys, name)
	}

	sort.Strings(keys)

	return keys
}

func objectPath(parent, name string) string {
	if parent == "" {
		return name
	}

	return parent + "." + name
}
