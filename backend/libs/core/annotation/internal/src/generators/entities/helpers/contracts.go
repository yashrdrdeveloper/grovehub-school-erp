package helpers

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"
)

type entitySpec struct {
	name             string
	fields           []fieldSpec
	builder, getter, setter bool
}

type fieldSpec struct {
	name, typ string
}

type sourceFile struct {
	name, pkg string
	imports []layout.ImportBinding
	entities []entitySpec
}

type methodIndex map[string]map[string]bool
