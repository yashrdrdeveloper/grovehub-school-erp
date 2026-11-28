// Package layout is a compatibility bridge for the annotation generator layout.
// Deprecated: maintained generators use the capability-private layout package.
package layout

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"

const CoreImportPath = layout.CoreImportPath

func ReadModulePath(root string) (string, error) { return layout.ReadModulePath(root) }
func ResolveImport(root, modulePath, importPath string) (string, error) { return layout.ResolveImport(root, modulePath, importPath) }
