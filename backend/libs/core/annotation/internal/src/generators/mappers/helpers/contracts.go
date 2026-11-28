package helpers

import (
	"grovehub/libs/core/annotation/internal/src/helpers/layout"
)

type mappingMethod struct {
	name, signature, output string
	fields                  []fieldMapping
}

type fieldMapping struct{ name, expression string }

type mappingGroup struct {
	receiver string
	imports  []layout.ImportBinding
	methods  []mappingMethod
}

// HasMethods reports whether discovery found executable mapping methods.
func HasMethods(group mappingGroup) bool {
	return len(group.methods) > 0
}
