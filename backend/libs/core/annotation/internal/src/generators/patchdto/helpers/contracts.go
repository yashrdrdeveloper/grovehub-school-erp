package helpers

const patchImport = "grovehub/libs/core/annotation"

type patchField struct {
	name string
	jsonName string
}

type patchType struct {
	name string
	fields []patchField
}
