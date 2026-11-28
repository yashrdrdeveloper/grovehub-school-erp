package helpers

import (
	"strconv"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/diff/internal/src/contracts"
)

func walkArray(path string, old, next []any, changes *[]contracts.Change) {
	for index := 0; index < max(len(old), len(next)); index++ {
		child := path + "[" + strconv.Itoa(index) + "]"
		walk(child, arrayValue(old, index), arrayValue(next, index), changes)
	}
}

func arrayValue(values []any, index int) any {
	if index >= len(values) {
		return nil
	}

	return values[index]
}
