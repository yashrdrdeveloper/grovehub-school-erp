package helpers

const patchImport = "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation"

type patchField struct {
	name string
	jsonName string
}

type patchType struct {
	name string
	fields []patchField
}
